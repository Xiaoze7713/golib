package memcache

import "sync"

type MemCache struct {
	rwLock sync.RWMutex
	kvMap map[string]interface{}
}

func NewMemCache() (m *MemCache) {
	m = &MemCache{
		rwLock: sync.RWMutex{},
		kvMap: map[string]interface{}{},
	}
	return
}

func (mc *MemCache) Get(key string) (value interface{}, ok bool) {
	mc.rwLock.RLock()
	value, ok = mc.kvMap[key]
	mc.rwLock.RUnlock()
	if !ok {
		return
	}
	return
}

func (mc *MemCache) Set(key string, value interface{}) (ok bool) {
	mc.rwLock.Lock()
	mc.kvMap[key] = value
	mc.rwLock.Unlock()
	return true
}

func (mc *MemCache) Del(key string) {
	mc.rwLock.Lock()
	delete(mc.kvMap, key)
	mc.rwLock.Unlock()
}