package encode

import (
	"chatbot/utils/constant"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/net/html/charset"
	"io"
	"strings"
)

func utf8Convert(state *lua.LState) int {
	lData := state.Get(constant.Param1)
	if lData.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: data must be a string")
		return 0
	}
	lContentType := state.Get(constant.Param2)
	if lContentType.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: content-type must be a string")
		return 0
	}

	data := strings.NewReader(string(lData.(lua.LString)))
	contentType := string(lContentType.(lua.LString))
	utf8Reader, err := charset.NewReader(data, contentType)
	if err != nil {
		state.Push(lua.LNil)
		state.Push(lua.LString(err.Error()))
		return 1
	}

	bd, err := io.ReadAll(utf8Reader)
	if err != nil {
		state.Push(lua.LNil)
		state.Push(lua.LString(err.Error()))
	}
	state.Push(lua.LString(bd))
	return 1
}
