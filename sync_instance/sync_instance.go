package sync_instance

import (
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
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
	newF     func(config *C) (*T, error)
	config   *C
}

func NewSyncMgr[T SyncHandler, C ConfigAny](newT func(config *C) (*T, error), config *C, interval time.Duration) (*SyncManager[T, C], error) {
	t, err := newT(config)
	return &SyncManager[T, C]{
		lock:     sync.RWMutex{},
		handler:  t,
		interval: interval,
		newF:     newT,
		config:   config,
	}, err
}

func (m *SyncManager[T, C]) LoopUpdate() {
	for {
		time.Sleep(m.interval)
		m.update()
	}
}

func (m *SyncManager[T, C]) update() {
	h, err := m.newF(m.config)
	if err != nil {
		xlog.Error(err)
		return
	}
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
