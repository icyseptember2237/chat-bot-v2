package encode

import (
	"chatbot/utils/engine_pool"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
)

const moduleName = "encode"

var moduleMethods = map[string]lua.LGFunction{
	"base64":      base64Encode,
	"utf8Convert": utf8Convert,
}

func init() {
	engine_pool.RegisterModule(func(i engine.Engine) {
		i.(*engine.LuaEngine).GetVM().PreloadModule(moduleName, func(state *lua.LState) int {
			module := state.SetFuncs(state.NewTable(), moduleMethods)
			state.SetField(module, "name", lua.LString(moduleName))
			state.Push(module)

			return 1
		})
	})
}
