package utils

import (
	"sync"
	"time"
)

type item[T any] struct {
	value    T
	initTime int64
}

type TTLMap[T any] struct {
	m map[string]*item[T]
	l sync.Mutex
}

func NewTTLMap[T any](ln int, maxTTL int64) (m *TTLMap[T]) {
	m = &TTLMap[T]{m: make(map[string]*item[T], ln)}
	go func() {
		for now := range time.Tick(time.Second) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix()-v.initTime > maxTTL {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return
}

func (m *TTLMap[T]) Len() int {
	return len(m.m)
}

func (m *TTLMap[T]) Put(k string, v T) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &item[T]{value: v}
		m.m[k] = it
	}
	it.initTime = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLMap[T]) Get(k string) (v T) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		v = it.value
		//it.lastAccess = time.Now().Unix()
	}
	m.l.Unlock()
	return
}
