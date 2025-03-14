package bot_api

import (
	"chatbot/utils/bot_api"
	"chatbot/utils/constant"
	"chatbot/utils/engine_pool"
	"chatbot/utils/luatool"
	"encoding/json"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
)

const moduleName = "bot_api"

var moduleMethods = map[string]lua.LGFunction{
	"get_group_member_list": getGroupMemberList,
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

func getGroupMemberList(state *lua.LState) int {
	lGroupId := state.Get(constant.Param1)
	if lGroupId.Type() != lua.LTNumber {
		state.ArgError(constant.Param1, "type error: groupId must be number")
		return 0
	}

	groupId := int64(lGroupId.(lua.LNumber))
	res, err := bot_api.GetGroupMemberList(groupId)
	if err != nil {
		state.Push(lua.LNil)
		return 1
	}

	js, _ := json.Marshal(res)
	var list []interface{}
	if err := json.Unmarshal(js, &list); err != nil {
		state.Push(lua.LNil)
		return 1
	}

	lTable := luatool.ConvertToTable(state, list)
	state.Push(lTable)
	return 1
}
