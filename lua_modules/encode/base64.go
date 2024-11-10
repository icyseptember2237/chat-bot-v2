package encode

import (
	"chatbot/utils/constant"
	"encoding/base64"
	lua "github.com/yuin/gopher-lua"
)

func base64Encode(state *lua.LState) int {
	lData := state.Get(constant.Param1)
	if lData.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: data must be string")
		return 0
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(lData.(lua.LString)))

	state.Push(lua.LString(encoded))
	return 1
}
