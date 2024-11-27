package go_query

import (
	"chatbot/utils/constant"
	"chatbot/utils/engine_pool"
	"github.com/PuerkitoBio/goquery"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

const moduleName = "goquery"

var moduleMethods = map[string]lua.LGFunction{
	"new": newDoc,
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

func newDoc(state *lua.LState) int {
	lData := state.Get(constant.Param1)
	if lData.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: data must be a string")
		return 0
	}
	data := strings.NewReader(string(lData.(lua.LString)))

	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		state.Push(lua.LNil)
		return 0
	}

	ud := state.NewUserData()
	ud.Value = doc

	state.SetMetatable(ud, state.GetTypeMetatable(metaName))
	state.Push(ud)
	return 1
}
