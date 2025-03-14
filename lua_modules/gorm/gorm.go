package gorm

import (
	"chatbot/utils/constant"
	"chatbot/utils/engine_pool"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
)

const moduleName = "gorm"

var moduleMethods = map[string]lua.LGFunction{
	"new": newClient,
}

func init() {
	engine_pool.RegisterModule(func(i engine.Engine) {
		i.(*engine.LuaEngine).GetVM().PreloadModule(moduleName, func(state *lua.LState) int {
			mod := state.SetFuncs(state.NewTable(), moduleMethods)
			state.SetField(mod, "name", lua.LString(moduleName))
			state.Push(mod)

			meta := state.NewTypeMetatable(metaName)
			state.SetField(meta, "__index", state.SetFuncs(state.NewTable(), clientExports))
			return 1
		})
	})
}

func newClient(state *lua.LState) int {
	lKind := state.Get(constant.Param1)
	if lKind.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: type must be a string")
	}
	lName := state.Get(constant.Param2)
	if lName.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: name must be a string")
	}

	cli := &client{
		name: lName.(lua.LString).String(),
		kind: lKind.(lua.LString).String(),
	}

	if ok := cli.init(); !ok {
		state.Error(lua.LString("client init failed"), 1)
		return 0
	}

	ud := state.NewUserData()
	ud.Value = cli
	state.SetMetatable(ud, state.GetTypeMetatable(moduleName))
	state.Push(ud)
	return 1
}
