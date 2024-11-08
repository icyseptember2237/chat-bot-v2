package mongo

import (
	"chatbot/logger"
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

func deleteOne(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: filter must be table")
		return 0
	}

	filter := gluamapper.ToGoValue(lFilter, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	filter = luatool.ConvertLuaData(filter)

	bsm, err := convertToBson(filter)
	if err != nil {
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	res, err := conn.col.DeleteOne(context.Background(), bsm)
	if err != nil {
		logger.Errorf(context.Background(), "col.DeleteOne error %s", err.Error())
		state.Push(lua.LNumber(0))
		return 1
	}

	state.Push(lua.LNumber(res.DeletedCount))
	return 1
}

func deleteMany(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: filter must be table")
		return 0
	}

	filter := gluamapper.ToGoValue(lFilter, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	filter = luatool.ConvertLuaData(filter)

	bsm, err := convertToBson(filter)
	if err != nil {
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	res, err := conn.col.DeleteMany(context.Background(), bsm)
	if err != nil {
		logger.Errorf(context.Background(), "col.DeleteMany error %s", err.Error())
		state.Push(lua.LNumber(0))
		return 0
	}

	state.Push(lua.LNumber(res.DeletedCount))
	return 1
}
