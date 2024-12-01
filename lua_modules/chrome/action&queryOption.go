package chrome

import (
	"chatbot/utils/constant"
	"github.com/chromedp/chromedp"
	lua "github.com/yuin/gopher-lua"
	"time"
)

var queryOptions = map[string]chromedp.QueryOption{
	"byQuery":    chromedp.ByQuery,
	"byID":       chromedp.ByID,
	"byJSPath":   chromedp.ByJSPath,
	"byQueryAll": chromedp.ByQueryAll,
	"bySearch":   chromedp.BySearch,
	"byNodeID":   chromedp.ByNodeID,
}

func navigate(state *lua.LState) int {
	lUrl := state.Get(constant.Param1)
	if lUrl.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: url must be a string")
		return 0
	}
	url := string(lUrl.(lua.LString))

	action := chromedp.Navigate(url)

	ud := state.NewUserData()
	ud.Value = action
	state.Push(ud)
	return 1
}

func navigateBack(state *lua.LState) int {
	action := chromedp.NavigateBack()

	ud := state.NewUserData()
	ud.Value = action
	state.Push(ud)
	return 1
}

func navigateForward(state *lua.LState) int {
	action := chromedp.NavigateForward()

	ud := state.NewUserData()
	ud.Value = action
	state.Push(ud)
	return 1
}

func waitVisible(state *lua.LState) int {
	lSelect := state.Get(constant.Param1)
	if lSelect.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: url must be a string")
		return 0
	}
	lSelector := state.Get(constant.Param2)
	if lSelector.Type() != lua.LTUserData {
		state.ArgError(constant.Param2, "type error: selector type error")
		return 0
	}
	sel := string(lSelect.(lua.LString))
	selector, ok := lSelector.(*lua.LUserData).Value.(chromedp.QueryOption)
	if !ok {
		state.ArgError(constant.Param2, "selector assert error")
		return 0
	}

	action := chromedp.WaitVisible(sel, selector)

	ud := state.NewUserData()
	ud.Value = action
	state.Push(ud)
	return 1
}

func click(state *lua.LState) int {
	lSelect := state.Get(constant.Param1)
	if lSelect.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: url must be a string")
		return 0
	}
	lSelector := state.Get(constant.Param2)
	if lSelector.Type() != lua.LTUserData {
		state.ArgError(constant.Param2, "type error: selector type error")
		return 0
	}
	sel := string(lSelect.(lua.LString))
	selector, ok := lSelector.(*lua.LUserData).Value.(chromedp.QueryOption)
	if !ok {
		state.ArgError(constant.Param2, "selector assert error")
		return 0
	}

	action := chromedp.Click(sel, selector)

	ud := state.NewUserData()
	ud.Value = action
	state.Push(ud)
	return 1
}

func sleep(state *lua.LState) int {
	lDuration := state.Get(constant.Param1)
	if lDuration.Type() != lua.LTNumber {
		state.ArgError(constant.Param1, "type error: duration must be a number")
		return 0
	}
	duration := time.Duration(lDuration.(lua.LNumber))

	action := chromedp.Sleep(duration * time.Second)

	ud := state.NewUserData()
	ud.Value = action
	state.Push(ud)
	return 1
}

func setValue(state *lua.LState) int {
	lSelect := state.Get(constant.Param1)
	if lSelect.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: url must be a string")
		return 0
	}
	lValue := state.Get(constant.Param2)
	if lValue.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: value must be a string")
		return 0
	}
	lSelector := state.Get(constant.Param3)
	if lSelector.Type() != lua.LTUserData {
		state.ArgError(constant.Param3, "type error: selector type error")
		return 0
	}

	sel := string(lSelect.(lua.LString))
	value := string(lValue.(lua.LString))
	selector, ok := lSelector.(*lua.LUserData).Value.(chromedp.QueryOption)
	if !ok {
		state.ArgError(constant.Param2, "selector assert error")
		return 0
	}

	action := chromedp.SetValue(sel, value, selector)

	ud := state.NewUserData()
	ud.Value = action
	state.Push(ud)
	return 1
}

func value(state *lua.LState) int {
	lSelect := state.Get(constant.Param1)
	if lSelect.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: select must be a string")
	}
	lSelector := state.Get(constant.Param2)
	if lSelector.Type() != lua.LTUserData {
		state.ArgError(constant.Param2, "type error: selector type error")
		return 0
	}

	sel := string(lSelect.(lua.LString))
	selector, ok := lSelector.(*lua.LUserData).Value.(chromedp.QueryOption)
	if !ok {
		state.ArgError(constant.Param2, "selector assert error")
		return 0
	}

	var v string
	action := chromedp.Value(sel, &v, selector)

	udRes := state.NewUserData()
	udRes.Value = &res{
		data: &v,
	}
	state.SetMetatable(udRes, state.GetTypeMetatable(pointerMeta))
	state.Push(udRes)

	udAction := state.NewUserData()
	udAction.Value = action
	state.Push(udAction)
	return 2
}

func attributeValue(state *lua.LState) int {
	lSelect := state.Get(constant.Param1)
	if lSelect.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: select must be a string")
	}
	lName := state.Get(constant.Param2)
	if lName.Type() != lua.LTString {
		state.ArgError(constant.Param2, "type error: name must be a string")
		return 0
	}
	lSelector := state.Get(constant.Param3)
	if lSelector.Type() != lua.LTUserData {
		state.ArgError(constant.Param2, "type error: selector type error")
		return 0
	}
	sel := string(lSelect.(lua.LString))
	name := string(lName.(lua.LString))
	selector, ok := lSelector.(*lua.LUserData).Value.(chromedp.QueryOption)
	if !ok {
		state.ArgError(constant.Param2, "selector assert error")
		return 0
	}

	var (
		v       string
		success bool
	)
	action := chromedp.AttributeValue(sel, name, &v, &success, selector)

	udRes := state.NewUserData()
	udRes.Value = &res{
		data:    &v,
		success: &success,
	}
	state.SetMetatable(udRes, state.GetTypeMetatable(pointerMeta))
	state.Push(udRes)

	udAction := state.NewUserData()
	udAction.Value = action
	state.Push(udAction)
	return 2
}

func title(state *lua.LState) int {
	var t string
	result := &res{
		data: &t,
	}

	action := chromedp.Title(&t)

	udRes := state.NewUserData()
	udRes.Value = result
	state.SetMetatable(udRes, state.GetTypeMetatable(pointerMeta))
	state.Push(udRes)

	udAction := state.NewUserData()
	udAction.Value = action
	state.Push(udAction)
	return 2
}

func evaluate(state *lua.LState) int {
	lExpression := state.Get(constant.Param1)
	if lExpression.Type() != lua.LTString {
		state.ArgError(constant.Param1, "type error: expression must be a string")
		return 0
	}
	expression := string(lExpression.(lua.LString))

	var result interface{}
	action := chromedp.Evaluate(expression, &result)

	udRes := state.NewUserData()
	udRes.Value = &res{
		data: &result,
	}
	state.SetMetatable(udRes, state.GetTypeMetatable(pointerMeta))
	state.Push(udRes)

	udAction := state.NewUserData()
	udAction.Value = action
	state.Push(udAction)

	return 2
}
