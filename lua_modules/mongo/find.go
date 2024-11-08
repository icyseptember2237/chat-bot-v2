package mongo

import (
	"chatbot/logger"
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	sortTableKey  = "sort"
	limitTableKey = "limit"
	skipTableKey  = "skip"
)

func findOne(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: filter must be table")
		return 0
	}

	filter := gluamapper.ToGoValue(lFilter, gluamapper.Option{NameFunc: func(s string) string { return s }})
	filter = luatool.ConvertLuaData(filter)

	bsm, err := convertToBson(filter)
	if err != nil {
		fmt.Println(err)
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	var res bson.Raw
	err = conn.col.FindOne(context.Background(), bsm).Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) == false {
			logger.Errorf(context.Background(), "col.FindOne error %s", err.Error())
		}
		state.Push(lua.LNil)
		return 1
	}

	result, err := convertFromBsonRaw(res)
	if err != nil {
		logger.Errorf(context.Background(), "convertFromBsonRaw error %s", err.Error())
		state.Push(lua.LNil)
		return 1
	}
	dst := luatool.ConvertToTable(state, result)
	state.Push(dst)
	return 1
}

func find(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: filter must be table")
		return 0
	}

	filter := gluamapper.ToGoValue(lFilter, gluamapper.Option{NameFunc: func(s string) string { return s }})
	filter = luatool.ConvertLuaData(filter)

	bsm, err := convertToBson(filter)
	if err != nil {
		fmt.Println(err)
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	lOption := state.Get(constant.Param3)
	var lSort, lLimit, lSkip lua.LValue
	if lOption.Type() == lua.LTTable {
		option := lOption.(*lua.LTable)
		lSort = option.RawGet(lua.LString(sortTableKey))
		lLimit = option.RawGet(lua.LString(limitTableKey))
		lSkip = option.RawGet(lua.LString(skipTableKey))
	}

	opts := options.Find()

	if lSort != nil && lSort.Type() == lua.LTTable {
		sort := gluamapper.ToGoValue(lSort, gluamapper.Option{NameFunc: func(s string) string { return s }})
		sort = luatool.ConvertLuaData(sort)

		stm, err := convertToBson(sort)
		if err == nil {
			opts.SetSort(stm)
		}
	}

	if lLimit != nil && lLimit.Type() == lua.LTNumber {
		limit := int64(lLimit.(lua.LNumber))
		opts.SetLimit(limit)
	}

	if lSkip != nil && lSkip.Type() == lua.LTNumber {
		skip := int64(lSkip.(lua.LNumber))
		opts.SetSkip(skip)
	}

	cursor, err := conn.col.Find(context.Background(), bsm, opts)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			logger.Errorf(context.Background(), "mongo.find %s error : %v", conn.name, err.Error())
		}
		state.Push(lua.LNil)
		return 1
	}

	ud := state.NewUserData()
	ud.Value = cursor

	state.SetMetatable(ud, state.GetTypeMetatable(cursorMetaName))
	state.Push(ud)
	return 1
}
