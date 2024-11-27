package go_query

import (
	"chatbot/utils/constant"
	"github.com/PuerkitoBio/goquery"
	lua "github.com/yuin/gopher-lua"
)

const metaName = "goquery{meta}"

var clientExports = map[string]lua.LGFunction{
	"find": find,
}

func find(state *lua.LState) int {
	doc := checkDoc(state)
	if doc == nil {
		return 0
	}

	lSelection := state.Get(constant.Param2)
	if lSelection.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type err: selection must be a string")
		return 0
	}
	selection := string(lSelection.(lua.LString))

	table := state.NewTable()
	doc.Find(selection).Each(func(i int, s *goquery.Selection) {
		content, err := s.Html()
		if err != nil {
			return
		}
		if content != "" {
			table.Append(lua.LString(content))
		}
	})

	state.Push(table)
	return 1
}

func checkDoc(state *lua.LState) *goquery.Document {
	ud := state.Get(constant.Param1)
	if ud.Type() != lua.LTUserData {
		state.ArgError(constant.Param1, "doc expected")
		return nil
	}

	if m, ok := ud.(*lua.LUserData).Value.(*goquery.Document); ok {
		return m
	}
	state.ArgError(constant.Param1, "msg empty")
	return nil
}
