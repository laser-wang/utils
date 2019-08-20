package utils

import (
	"sync"
)

type BeeMap struct {
	lock *sync.RWMutex
	BM   map[interface{}]interface{}
}

func NewBeeMap() *BeeMap {
	return &BeeMap{
		lock: new(sync.RWMutex),
		BM:   make(map[interface{}]interface{}),
	}
}

//Get from maps return the k's value
func (m *BeeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.BM[k]; ok {
		return val
	}
	return nil
}

// Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *BeeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	//	if val, ok := m.BM[k]; !ok {
	//		m.BM[k] = v
	//	} else if val != v {
	//		m.BM[k] = v
	//	} else {
	//		return false
	//	}
	m.BM[k] = v
	return true
}

// Returns true if k is exist in the map.
func (m *BeeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.BM[k]; !ok {
		return false
	}
	return true
}

func (m *BeeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.BM, k)
}
