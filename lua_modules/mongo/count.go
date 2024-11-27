package mongo

import (
	"chatbot/logger"
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	lua "github.com/yuin/gopher-lua"
)

func count(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type err: filter must be of type table")
		return 0
	}
	filter := ToGoValue(lFilter, func(s string) string {
		return s
	})
	filter = luatool.ConvertLuaData(filter)

	bd, err := convertToBson(filter)
	if err != nil {
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	count, err := conn.col.CountDocuments(context.Background(), bd)
	if err != nil {
		logger.Errorf(context.Background(), "mongo.count %s error : %v", conn.name, err.Error())
		state.Push(lua.LNumber(0))
		return 1
	}

	state.Push(lua.LNumber(count))
	return 1
}
