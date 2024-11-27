package mongo

import (
	"chatbot/utils/constant"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func oId(state *lua.LState) int {
	lOid := state.Get(constant.Param1)
	if lOid.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: oid must be a string")
		return 0
	}

	oid := string(lOid.(lua.LString))

	objectId, err := primitive.ObjectIDFromHex(oid)
	if err != nil {
		state.Error(lua.LString(err.Error()), 1)
		return 0
	}

	ud := state.NewUserData()
	ud.Value = objectId
	state.Push(ud)
	return 1
}
