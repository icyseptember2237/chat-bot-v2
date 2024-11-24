package task

import (
	"chatbot/logger"
	"chatbot/utils/constant"
	"chatbot/utils/engine_pool"
	"context"
	"github.com/icyseptember2237/engine"
	lua "github.com/yuin/gopher-lua"
	"time"
)

const moduleName = "task"

var moduleMethods = map[string]lua.LGFunction{
	"new": newTask,
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

func newTask(state *lua.LState) int {
	lFunc := state.Get(constant.Param1)
	if lFunc.Type() != lua.LTFunction {
		state.ArgError(constant.Param1, "type error: func must be function")
		return 0
	}
	lArgs := state.Get(constant.Param2)
	if lArgs.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: args must be table")
		return 0
	}

	lDelay := state.Get(constant.Param3)
	if lDelay.Type() != lua.LTNumber {
		state.ArgError(constant.Param3, "type error: delay must be number")
		return 0
	}

	go func() {
		function := lFunc.(*lua.LFunction)
		args := make([]lua.LValue, 0)
		lArgs.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
			args = append(args, value2)
		})
		delay := lDelay.(lua.LNumber)

		eng := engine_pool.GetDefaultEnginePool().GetRawEngine()
		vm := eng.(*engine.LuaEngine).GetVM()

		defer func() {
			vm.SetContext(context.Background())
			engine_pool.GetDefaultEnginePool().PutEngine(eng)
		}()

		t := time.After(time.Second * time.Duration(delay))
		<-t

		if err := vm.CallByParam(lua.P{
			Fn:      function,
			NRet:    0,
			Protect: true,
			Handler: nil,
		}, args...); err != nil {
			logger.Infof(context.Background(), "task %s run error: %v", function.String(), err)
		}
	}()

	return 0
}
