package job

import "sync"

var defaultManager *Manager

func init() {
	defaultManager = &Manager{}
}

func GetManager() *Manager {
	return defaultManager
}

type Manager struct {
	jobs sync.Map
}

func (m *Manager) Add(job *Job) bool {
	if _, ok := m.jobs.Load(job.conf.Name); ok {
		return false
	}

	m.jobs.Store(job.conf.Name, job)
	return true
}

func (m *Manager) Get(name string) *Job {
	if job, ok := m.jobs.Load(name); ok {
		return job.(*Job)
	} else {
		return nil
	}
}

func (m *Manager) Del(name string) {
	if job, ok := m.jobs.Load(name); ok {
		if job.(*Job).Running() {
			job.(*Job).Stop()
		}
		m.jobs.Delete(name)
	}
}

func (m *Manager) StartAll() {
	m.jobs.Range(func(key, value interface{}) bool {
		go func(job *Job) {
			job.Start()
		}(value.(*Job))
		return true
	})
}

func (m *Manager) StopAll() {
	m.jobs.Range(func(key, value interface{}) bool {
		go func(job *Job) {
			job.Stop()
		}(value.(*Job))
		return true
	})
}

func (m *Manager) GetAllNames() []string {
	names := make([]string, 0)
	m.jobs.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}
