package lua_msg

import (
	"chatbot/config"
	"chatbot/msg"
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"encoding/json"
	lua "github.com/yuin/gopher-lua"
)

const metaName = "group_msg{meta}"

var clientExports = map[string]lua.LGFunction{
	"textSegment":   textSegment,
	"atSegment":     atSegment,
	"replySegment":  replySegment,
	"imageSegment":  imageSegment,
	"fileSegment":   fileSegment,
	"recordSegment": recordSegment,
	"toJsonString":  toJsonString,
	"send":          send,
}

func textSegment(state *lua.LState) int {
	m := checkMsg(state)
	if m == nil {
		return 0
	}

	lText := state.Get(constant.Param2)
	if lText.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: text must be string")
		return 1
	}

	text := lText.(lua.LString).String()

	m.AppendSegment(msg.NewTextSegment(text))
	return 0
}

func atSegment(state *lua.LState) int {
	m := checkMsg(state)
	if m == nil {
		return 0
	}

	lAt := state.Get(constant.Param2)
	if lAt.Type() != lua.LTNumber {
		state.ArgError(constant.Param2, "type error: at must be number")
		return 1
	}
	at := int64(lAt.(lua.LNumber))

	m.AppendSegment(msg.NewAtSegment(at))
	return 0
}

func replySegment(state *lua.LState) int {
	m := checkMsg(state)
	if m == nil {
		return 0
	}

	lMsgId := state.Get(constant.Param2)
	if lMsgId.Type() != lua.LTNumber {
		state.ArgError(constant.Param2, "type error: reply must be number")
		return 1
	}
	msgId := int64(lMsgId.(lua.LNumber))

	m.AppendSegment(msg.NewReplySegment(msgId))
	return 0
}

func imageSegment(state *lua.LState) int {
	m := checkMsg(state)
	if m == nil {
		return 0
	}

	lImg := state.Get(constant.Param2)
	if lImg.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: image must be string")
		return 1
	}
	img := lImg.(lua.LString).String()

	m.AppendSegment(msg.NewImageSegment(img))
	return 0
}

func fileSegment(state *lua.LState) int {
	m := checkMsg(state)
	if m == nil {
		return 0
	}

	lFile := state.Get(constant.Param2)
	if lFile.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: file must be string")
		return 1
	}
	file := lFile.(lua.LString).String()

	lName := state.Get(constant.Param3)
	if lName.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: name must be string")
		return 1
	}
	name := lName.(lua.LString).String()

	m.AppendSegment(msg.NewFileSegment(file, name))
	return 0
}

func recordSegment(state *lua.LState) int {
	m := checkMsg(state)
	if m == nil {
		return 0
	}

	lFile := state.Get(constant.Param2)
	if lFile.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: file must be string")
		return 1
	}
	file := lFile.(lua.LString).String()

	m.AppendSegment(msg.NewRecordSegment(file))
	return 0
}

func toJsonString(state *lua.LState) int {
	m := checkMsg(state)
	if m == nil {
		return 0
	}

	j, err := json.Marshal(m)
	if err != nil {
		state.Push(lua.LNil)
		return 1
	}
	state.Push(lua.LString(j))
	return 1
}

func send(state *lua.LState) int {
	m := checkMsg(state)
	if m == nil {
		return 0
	}
	conf := config.Get().Server
	if res, err := m.Send(conf.BotAddr, conf.BotToken); err == nil {
		lTable := luatool.ConvertToTable(state, res)
		state.Push(lua.LTrue)
		state.Push(lTable)
		return 1
	} else {
		state.Push(lua.LFalse)
		state.Push(lua.LString(err.Error()))
		return 0
	}
}

func checkMsg(state *lua.LState) *msg.GroupMessage {
	ud := state.Get(constant.Param1)
	if ud.Type() != lua.LTUserData {
		state.ArgError(constant.Param1, "client expected")
		return nil
	}

	if m, ok := ud.(*lua.LUserData).Value.(*msg.GroupMessage); ok {
		return m
	}
	state.ArgError(constant.Param1, "msg empty")
	return nil
}
