package engine_pool

var defaultEP *EnginePool

func init() {
	defaultEP = NewEnginePool()
}

func GetDefaultEnginePool() *EnginePool {
	return defaultEP
}
