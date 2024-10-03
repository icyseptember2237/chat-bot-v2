package worker

import "sync"

func init() {
	defaultManager = &Manager{}
}

func GetManager() *Manager {
	return defaultManager
}

var defaultManager *Manager

type Manager struct {
	workers sync.Map
}

func (m *Manager) Add(worker *Worker) bool {
	if _, ok := m.workers.Load(worker.conf.Name); ok {
		return false
	}

	m.workers.Store(worker.conf.Name, worker)
	return true
}

func (m *Manager) Get(name string) *Worker {
	if worker, ok := m.workers.Load(name); ok {
		return worker.(*Worker)
	} else {
		return nil
	}
}

func (m *Manager) Del(name string) {
	if worker, ok := m.workers.Load(name); ok {
		if worker.(*Worker).Running() {
			worker.(*Worker).Stop()
		}
		m.workers.Delete(name)
	}
}

func (m *Manager) StartAll() {
	m.workers.Range(func(key, value interface{}) bool {
		go func(worker *Worker) {
			worker.Start()
		}(value.(*Worker))
		return true
	})
}

func (m *Manager) StopAll() {
	m.workers.Range(func(key, value interface{}) bool {
		go func(worker *Worker) {
			worker.Stop()
		}(value.(*Worker))
		return true
	})
}

func (m *Manager) GetAllNames() []string {
	names := make([]string, 0)
	m.workers.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}
