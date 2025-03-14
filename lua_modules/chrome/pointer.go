package chrome

import (
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	lua "github.com/yuin/gopher-lua"
	"reflect"
)

const pointerMeta = "pointer{meta}"

var pointerExports = map[string]lua.LGFunction{
	"getData": get,
	"success": success,
	"set":     set,
}

type res struct {
	data    interface{}
	success *bool
}

func get(state *lua.LState) int {
	ptr := checkPointer(state)
	if ptr == nil {
		return 0
	}

	switch v := ptr.data.(type) {
	case *string:
		state.Push(lua.LString(*v))
		break
	case *interface{}:
		val := reflect.ValueOf(v).Elem()
		if !val.IsValid() {
			state.Push(lua.LNil)
			break
		}
		valElem := val.Elem()
		if !valElem.IsValid() {
			state.Push(lua.LNil)
			break
		}
		state.Push(luatool.ConvertToTable(state, valElem.Interface()))
		break
	default:
		break
	}

	return 1
}

func success(state *lua.LState) int {
	ptr := checkPointer(state)
	if ptr == nil {
		return 0
	}

	state.Push(lua.LBool(*ptr.success))
	return 1
}

func set(state *lua.LState) int {
	return 0
}

func checkPointer(state *lua.LState) *res {
	lP := state.Get(constant.Param1)
	if lP.Type() != lua.LTUserData {
		state.ArgError(constant.Param1, "userdata expected")
		return nil
	}
	p := lP.(*lua.LUserData).Value.(*res)
	return p
}
