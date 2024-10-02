package engine_pool

import "github.com/icyseptember2237/engine"

type EnginePool struct {
	egs *engine.EnginePool
}

func NewEnginePool() *EnginePool {
	return &EnginePool{
		egs: engine.InitEnginePool(engine.TypeEngineLua),
	}
}

func (ep *EnginePool) GetEngine(script string) (engine.Engine, error, bool) {
	eng := ep.egs.Get()
	if !eng.IsReady() {
		for _, fn := range moduleFun {
			fn(eng)
		}

		err := eng.ParseFile(script)
		if err != nil {
			return nil, err, false
		}

		eng.SetReady()
		return eng, nil, true
	}
	return eng, nil, true
}

func (ep *EnginePool) GetRawEngine() engine.Engine {
	eng := ep.egs.Get()
	if !eng.IsReady() {
		for _, fn := range moduleFun {
			fn(eng)
		}

		eng.SetReady()
	}
	return eng
}

func (ep *EnginePool) PutEngine(e engine.Engine) {
	if e != nil {
		ep.egs.Put(e)
	}
}
