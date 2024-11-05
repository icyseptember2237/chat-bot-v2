package mongo

import (
	"chatbot/storage/mongo"
	"chatbot/utils/constant"
	"chatbot/utils/engine_pool"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
)

const moduleName = "mongo"

var moduleMethods = map[string]lua.LGFunction{
	"new": newConn,
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

func newConn(state *lua.LState) int {
	lName := state.Get(constant.Param1)
	if lName.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: name must be string")
		return 0
	}
	lBase := state.Get(constant.Param2)
	if lBase.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: base must be string")
		return 0
	}
	lCollection := state.Get(constant.Param3)
	if lCollection.Type() != lua.LTString {
		state.ArgError(constant.Param3, "type error: collection must be string")
		return 0
	}

	c := &Connection{
		name:       lName.(lua.LString).String(),
		database:   lBase.(lua.LString).String(),
		collection: lCollection.(lua.LString).String(),
	}

	cli := mongo.Get(c.name)
	if cli == nil {
		state.ArgError(constant.Param1, "param error: client is nil")
		return 0
	}

	c.col = cli.Database(c.database).Collection(c.collection)

	ud := state.NewUserData()
	ud.Value = c

	state.SetMetatable(ud, state.GetTypeMetatable(metaName))
	state.Push(ud)
	return 1
}
