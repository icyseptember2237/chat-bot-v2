package engine_pool

import (
	"github.com/icyseptember2237/engine"
)

var moduleFun []func(engine.Engine)

func RegisterModule(fn func(engine.Engine)) {
	moduleFun = append(moduleFun, fn)
}
