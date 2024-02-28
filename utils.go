package dbg

import "sync"

func getOk[K any, V any](m *sync.Map, key K) (V, bool) {
	val, ok := m.Load(key)
	return val.(V), ok
}
