package mongo

import (
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func aggregate(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: filter must be table")
		return 0
	}

	filter := ToGoValue(lFilter, func(s string) string { return s })
	filter = luatool.ConvertLuaData(filter)

	bs, err := convertToBson(filter)
	if err != nil {
		fmt.Println(err)
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	cursor, err := conn.col.Aggregate(context.Background(), bs)
	if err != nil {
		state.Push(lua.LNil)
		return 1
	}

	ud := state.NewUserData()
	ud.Value = cursor

	state.SetMetatable(ud, state.GetTypeMetatable(cursorMetaName))
	state.Push(ud)

	return 1
}
