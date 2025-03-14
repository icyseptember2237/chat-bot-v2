package gorm

import (
	"chatbot/logger"
	connection "chatbot/storage/gorm"
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	"database/sql"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
	"strconv"
)

const metaName = "gorm{meta}"

var clientExports = map[string]lua.LGFunction{
	"insert": insert,
	"update": update,
	"run":    run,
}

func insert(state *lua.LState) int {
	cli := checkClient(state)
	if cli == nil {
		return 0
	}

	lTable := state.Get(constant.Param2)
	if lTable.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type err: table must be of type string")
		return 0
	}
	lData := state.Get(constant.Param3)
	if lData.Type() != lua.LTTable {
		state.ArgError(constant.Param3, "type err: data must be of type table")
		return 0
	}

	data := gluamapper.ToGoValue(lData, gluamapper.Option{NameFunc: gluamapper.ToUpperCamelCase})
	data = luatool.ConvertLuaData(data)

	tx := cli.db.Table(lTable.(lua.LString).String()).Create(data)

	state.Push(lua.LNumber(tx.RowsAffected))
	return 1
}

func update(state *lua.LState) int {
	cli := checkClient(state)
	if cli == nil {
		return 0
	}

	lTable := state.Get(constant.Param2)
	if lTable.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type err: table must be of type string")
		return 0
	}
	lQuery := state.Get(constant.Param3)
	if lQuery.Type() != lua.LTString {
		state.ArgError(constant.Param3, "type err: query must be of type string")
		return 0
	}
	lData := state.Get(constant.Param4)
	if lData.Type() != lua.LTTable {
		state.ArgError(constant.Param4, "type err: data must be of type table")
		return 0
	}
	var limit int
	lLimit := state.Get(constant.Param5)
	if lLimit.Type() == lua.LTNil {
		limit = 0
	} else if lLimit.Type() != lua.LTNumber {
		state.ArgError(constant.Param5, "type err: limit must be of type number")
		return 0
	} else {
		var err error
		limit, err = strconv.Atoi(lLimit.(lua.LNumber).String())
		if err != nil {
			state.ArgError(constant.Param4, "type err: limit must be of type integer (limit)")
			return 0
		}
	}

	data := gluamapper.ToGoValue(lData, gluamapper.Option{NameFunc: gluamapper.ToUpperCamelCase})
	data = luatool.ConvertLuaData(data)

	var tx *gorm.DB
	if limit > 0 {
		tx = cli.db.Table(lTable.(lua.LString).String()).Where(lQuery.(lua.LString).String()).Updates(data).Limit(limit)
	} else {
		tx = cli.db.Table(lTable.(lua.LString).String()).Where(lQuery.(lua.LString).String()).Updates(data)
	}
	state.Push(lua.LNumber(tx.RowsAffected))
	return 1
}

func run(state *lua.LState) int {
	cli := checkClient(state)
	if cli == nil {
		return 0
	}

	lSql := state.Get(constant.Param2)
	if lSql.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type err: sql must be of type string")
		return 0
	}
	sqlStr := lSql.(lua.LString).String()

	lParams := state.Get(constant.Param3)
	var params []interface{}
	if lParams.Type() == lua.LTNil {
	} else if lParams.Type() == lua.LTTable {
		lParams.(*lua.LTable).ForEach(func(key lua.LValue, val lua.LValue) {
			if val.Type() == lua.LTString {
				params = append(params, val.(lua.LString).String())
			} else if val.Type() == lua.LTBool {
				params = append(params, val.(lua.LBool).String())
			} else if val.Type() == lua.LTNumber {
				params = append(params, val.(lua.LNumber).String())
			} else if val.Type() == lua.LTTable {
				var subParams []interface{}
				val.(*lua.LTable).ForEach(func(subKey lua.LValue, subVal lua.LValue) {
					if subVal.Type() == lua.LTString {
						subParams = append(subParams, subVal.(lua.LString).String())
					} else if subVal.Type() == lua.LTNumber {
						subParams = append(subParams, subVal.(lua.LNumber).String())
					}
				})
				params = append(params, subParams)
			}
		})
	} else {
		state.ArgError(constant.Param3, "type err: params must be of type table")
		return 0
	}

	var rows *sql.Rows
	var err error
	if len(params) > 0 {
		rows, err = cli.db.Raw(sqlStr, params...).Rows()
	} else {
		rows, err = cli.db.Raw(sqlStr).Rows()
	}
	if err != nil {
		logger.Errorf(context.Background(), "sql %s params %+v error %s", sqlStr, params, err.Error())
		state.Push(lua.LNil)
		return 1
	}

	lRes := state.NewTable()
	for rows.Next() {
		data := make(map[string]interface{})
		err = cli.db.ScanRows(rows, &data)
		dst := luatool.ConvertToTable(state, data)
		lRes.Append(dst)
	}

	state.Push(lRes)
	return 1
}

func checkClient(state *lua.LState) *client {
	ud := state.Get(constant.Param1)
	if ud.Type() != lua.LTUserData {
		state.ArgError(constant.Param1, "client expected")
		return nil
	}

	if cli, ok := ud.(*lua.LUserData).Value.(*client); ok {
		if cli.db == nil {
			state.ArgError(constant.Param1, "client no connection")
			return nil
		}

		if err := cli.ping(); err != nil {
			state.Error(lua.LString(err.Error()), 1)
			return nil
		}
		return cli
	}
	state.ArgError(constant.Param1, "client expected")
	return nil
}

type client struct {
	name string
	kind string
	db   *gorm.DB
}

func (c *client) ping() error {
	db, err := c.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (c *client) init() bool {
	var (
		key string
		ok  bool
	)
	switch c.kind {
	case "mysql":
		key = "mysql_" + c.name
	case "postgres":
		key = "postgres_" + c.name
	default:
		return false
	}

	c.db, ok = connection.Get(key)
	return ok
}
