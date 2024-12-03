package chrome

import (
	"chatbot/config"
	"chatbot/logger"
	"chatbot/utils/engine_pool"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
	"sync"
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
var mu sync.Mutex

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
}

func checkRemoteAllocator() bool {
	if allocator == nil {
		mu.Lock()
		defer mu.Unlock()
		if allocator == nil {
			allocator, _ = chromedp.NewRemoteAllocator(context.Background(), config.Get().Resources.Chromium)
			return allocator != nil
		}
	}
	return true
}

func newContext(state *lua.LState) int {
	if !checkRemoteAllocator() {
		logger.Warnf(context.Background(), "allocator init failed")
		state.Push(lua.LNil)
		return 1
	}

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
