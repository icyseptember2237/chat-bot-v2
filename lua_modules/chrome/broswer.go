package chrome

import (
	"chatbot/logger"
	"chatbot/utils/constant"
	"context"
	"github.com/chromedp/chromedp"
	lua "github.com/yuin/gopher-lua"
)

const browserMeta = "browser{meta}"

var browserExports = map[string]lua.LGFunction{
	"run":   run,
	"close": closeBrowser,
}

type browserContext struct {
	ctx  context.Context
	opts []chromedp.ExecAllocatorOption
}

func run(state *lua.LState) int {
	ctx := checkBrowser(state)
	if ctx == nil {
		return 0
	}

	lActions := state.Get(constant.Param2)
	if lActions.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: actions must be table")
		return 0
	}

	var actions []chromedp.Action
	lActions.(*lua.LTable).ForEach(func(k, v lua.LValue) {
		if action, ok := v.(*lua.LUserData).Value.(chromedp.Action); ok {
			actions = append(actions, action)
		} else {
			logger.Debugf(context.Background(), "action assert err: %+v", v)
		}
	})

	if err := chromedp.Run(ctx, actions...); err != nil {
		logger.Infof(context.Background(), "chromedp.Run error: %v", err)
		state.Push(lua.LFalse)
		return 1
	}
	// TODO html元素解析返回
	return 1
}

func closeBrowser(state *lua.LState) int {
	ctx := checkBrowser(state)
	if ctx == nil {
		return 0
	}

	if err := chromedp.Cancel(ctx); err != nil {
		logger.Infof(context.Background(), "browser close error: %s", err)
	}
	return 0
}

func checkBrowser(state *lua.LState) context.Context {
	lBrowser := state.Get(constant.Param1)
	if lBrowser.Type() != lua.LTUserData {
		state.ArgError(constant.Param1, "browser expected")
		return nil
	}
	if b, ok := lBrowser.(*lua.LUserData).Value.(browserContext); ok {
		return b.ctx
	}
	state.ArgError(constant.Param1, "browser expected")
	return nil
}
