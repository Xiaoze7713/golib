package sync_instance

import (
	"sync"
	"time"
)

type ConfigAny interface {
	any
}

type SyncHandler interface {
	any
	//Load() *SyncHandler
}

type SyncManager[T SyncHandler, C ConfigAny] struct {
	lock     sync.RWMutex
	handler  *T
	interval time.Duration
	newF     func(config *C) *T
	config   *C
}

func NewSyncMgr[T SyncHandler, C ConfigAny](newT func(config *C) *T, config *C, interval time.Duration) *SyncManager[T, C] {
	t := newT(config)
	return &SyncManager[T, C]{
		lock:     sync.RWMutex{},
		handler:  t,
		interval: interval,
		newF:     newT,
		config:   config,
	}
}

func (m *SyncManager[T, C]) LoopUpdate() {
	for {
		time.Sleep(m.interval)
		m.update()
	}
}

func (m *SyncManager[T, C]) update() {
	h := m.newF(m.config)
	m.lock.Lock()
	defer m.lock.Unlock()
	m.handler = h
}

func (m *SyncManager[T, C]) GetInstance() *T {
	m.lock.RLock()
	defer m.lock.RUnlock()
	h := m.handler
	return h
}
