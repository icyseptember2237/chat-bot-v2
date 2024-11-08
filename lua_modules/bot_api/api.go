package bot_api

import (
	"chatbot/utils/bot_api"
	"chatbot/utils/engine_pool"
	"github.com/icyseptember2237/engine"
)

const moduleName = "bot_api"

var moduleMethods = map[string]interface{}{
	"get_group_member_list": getGroupMemberList,
}

func init() {
	engine_pool.RegisterModule(func(i engine.Engine) {
		i.RegisterModule(moduleName, moduleMethods)
	})
}

func getGroupMemberList(groupId float64) ([]bot_api.MemberInfo, error) {
	return bot_api.GetGroupMemberList(int64(groupId))
}
