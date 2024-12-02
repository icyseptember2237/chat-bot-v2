package chrome

import (
	"chatbot/logger"
	"chatbot/utils/engine_pool"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
	"time"
)

const moduleName = "chromedp"

var moduleMethods = map[string]lua.LGFunction{
	"newContext": newContext,
	//"options":    options,
	// chromedp.Action
	"navigate":        navigate,
	"navigateBack":    navigateBack,
	"navigateForward": navigateForward,
	"waitVisible":     waitVisible,
	"click":           click,
	"setValue":        setValue,
	"value":           value,
	"attributeValue":  attributeValue,
	"title":           title,
	"evaluate":        evaluate,
	"sleep":           sleep,
}

var allocator context.Context

func init() {
	engine_pool.RegisterModule(func(i engine.Engine) {
		i.(*engine.LuaEngine).GetVM().PreloadModule(moduleName, func(state *lua.LState) int {
			module := state.SetFuncs(state.NewTable(), moduleMethods)
			for name, opt := range queryOptions {
				ud := state.NewUserData()
				ud.Value = opt
				state.SetField(module, name, ud)
			}

			state.SetField(module, "name", lua.LString(moduleName))
			state.Push(module)

			bMeta := state.NewTypeMetatable(browserMeta)
			state.SetField(bMeta, "__index", state.SetFuncs(state.NewTable(), browserExports))

			pMeta := state.NewTypeMetatable(pointerMeta)
			state.SetField(pMeta, "__index", state.SetFuncs(state.NewTable(), pointerExports))
			return 1
		})
	})

	option := append(chromedp.DefaultExecAllocatorOptions[:],
		//无头模式
		chromedp.Flag("headless", true),
		// 防止监测webdriver
		chromedp.Flag("enable-automation", false),
		//禁用 blink 特征
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		//忽略浏览器的风险提示（但好像并没什么用）
		chromedp.Flag("ignore-certificate-errors", true),
		//关闭浏览器声音（也没用）
		chromedp.Flag("mute-audio", false),
		//设置浏览器尺寸
		chromedp.WindowSize(1150, 1000),
	)
	allocator, _ = chromedp.NewExecAllocator(context.Background(), option...)
}

func newContext(state *lua.LState) int {

	ctx, _ := chromedp.NewContext(allocator, chromedp.WithLogf(func(s string, i ...interface{}) {
		logger.Infof(context.Background(), "chrome ctx run: %s", fmt.Sprintf(s, i...))
	}))

	go func() {
		select {
		case <-time.After(30 * time.Minute):
			if err := chromedp.Cancel(ctx); err != nil {
				logger.Infof(context.Background(), "browser close error: %s", err)
			}
			logger.Infof(context.Background(), "browser exited")
			return
		case <-ctx.Done():
			logger.Infof(context.Background(), "browser exited")
			return
		}
	}()

	ud := state.NewUserData()
	ud.Value = browserContext{
		ctx: ctx,
	}
	state.SetMetatable(ud, state.GetTypeMetatable(browserMeta))
	state.Push(ud)
	return 1
}

//func options(state *lua.LState) int {
//	lOptions := state.Get(constant.Param1)
//	if lOptions.Type() != lua.LTTable {
//		state.ArgError(constant.Param1, "type err: option must be a table")
//		return 0
//	}
//
//	opts := chromedp.DefaultExecAllocatorOptions[:]
//	lOptions.(*lua.LTable).ForEach(func(k, v lua.LValue) {
//		name := string(k.(lua.LString))
//		value := bool(v.(lua.LBool))
//		opts = append(opts, chromedp.Flag(name, value))
//	})
//
//	ud := state.NewUserData()
//	ud.Value = opts
//	state.Push(ud)
//	return 1
//}
