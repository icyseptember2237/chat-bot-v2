package mongo

import (
	"chatbot/logger"
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	upsertTableKey = "upsert"
)

func updateOne(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: filter must be table")
		return 0
	}
	lUpdate := state.Get(constant.Param3)
	if lUpdate.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: update must be table")
		return 0
	}
	lOptions := state.Get(constant.Param4)
	var lUpsert lua.LValue
	if lOptions.Type() == lua.LTTable {
		optTable := lOptions.(*lua.LTable)
		lUpsert = optTable.RawGet(lua.LString(upsertTableKey))
	}

	filter := gluamapper.ToGoValue(lFilter, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	filter = luatool.ConvertLuaData(filter)

	filterM, err := convertToBson(filter)
	if err != nil {
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	update := gluamapper.ToGoValue(lUpdate, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	update = luatool.ConvertLuaData(update)

	updateM, err := convertToBson(update)
	if err != nil {
		state.ArgError(constant.Param2, "invalid update format")
		return 0
	}

	opts := options.Update()
	if lUpsert != nil && lUpsert.Type() != lua.LTBool {
		opts.SetUpsert(bool(lUpsert.(lua.LBool)))
	}

	res, err := conn.col.UpdateOne(context.Background(), filterM, updateM, opts)
	if err != nil {
		logger.Errorf(context.Background(), "col.UpdateByOne error %s", err.Error())
		state.Push(lua.LNumber(0))
		return 1
	}

	state.Push(lua.LNumber(res.ModifiedCount + res.UpsertedCount))
	return 1
}

func updateMany(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: filter must be table")
		return 0
	}
	lUpdate := state.Get(constant.Param3)
	if lUpdate.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: update must be table")
		return 0
	}
	lOptions := state.Get(constant.Param4)
	var lUpsert lua.LValue
	if lOptions.Type() == lua.LTTable {
		optTable := lOptions.(*lua.LTable)
		lUpsert = optTable.RawGet(lua.LString(upsertTableKey))
	}

	filter := gluamapper.ToGoValue(lFilter, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	filter = luatool.ConvertLuaData(filter)

	filterM, err := convertToBson(filter)
	if err != nil {
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	update := gluamapper.ToGoValue(lUpdate, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	update = luatool.ConvertLuaData(update)

	updateM, err := convertToBson(update)
	if err != nil {
		state.ArgError(constant.Param2, "invalid update format")
		return 0
	}

	opts := options.Update()
	if lUpsert != nil && lUpsert.Type() != lua.LTBool {
		opts.SetUpsert(bool(lUpsert.(lua.LBool)))
	}

	res, err := conn.col.UpdateMany(context.Background(), filterM, updateM, opts)
	if err != nil {
		logger.Errorf(context.Background(), "col.UpdateMany error %s", err.Error())
		state.Push(lua.LNumber(0))
		return 1
	}

	state.Push(lua.LNumber(res.ModifiedCount + res.UpsertedCount))
	return 1
}
