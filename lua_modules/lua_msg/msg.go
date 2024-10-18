package lua_msg

import (
	"chatbot/msg"
	"chatbot/utils/constant"
	"chatbot/utils/engine_pool"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
)

const moduleName = "group_msg"

var moduleMethods = map[string]lua.LGFunction{
	"new": newGroupMessage,
}

func init() {
	engine_pool.RegisterModule(func(i engine.Engine) {
		i.(*engine.LuaEngine).GetVM().PreloadModule(moduleName, func(state *lua.LState) int {
			module := state.SetFuncs(state.NewTable(), moduleMethods)
			state.SetField(module, "name", lua.LString(moduleName))
			state.Push(module)

			meta := state.NewTypeMetatable(metaName)
			state.SetField(meta, "__index", state.SetFuncs(state.NewTable(), clientExports))
			return 1
		})
	})
}

func newGroupMessage(state *lua.LState) int {
	lReceiver := state.Get(constant.Param1)
	if lReceiver.Type() != lua.LTNumber {
		state.ArgError(constant.Param1, "type error: receiver must be number")
		return 0
	}
	receiver := int64(lReceiver.(lua.LNumber))

	m := msg.NewGroupMessage(receiver)
	ud := state.NewUserData()
	ud.Value = m

	state.SetMetatable(ud, state.GetTypeMetatable(metaName))
	state.Push(ud)
	return 1
}
