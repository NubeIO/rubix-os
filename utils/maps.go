package utils
import (
	"encoding/json"
	"sync"
)


func NewMap() *Map {
	return &Map{
		data: make(map[interface{}]interface{}),
	}
}

func NewMapFrom(v interface{}) (m *Map) {
	switch v := v.(type) {
	case map[interface{}]interface{}:
		return &Map{
			data: v,
		}
	case *Map:
		m = v
	}
	return
}

type Map struct {
	mutex sync.RWMutex
	data map[interface{}]interface{}
}

func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

func (m *Map) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.data)
}


func (m *Map) Set(key, value interface{}) *Map {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.set(key, value)
	return m
}

func (m *Map) Get(key interface{}) interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.getOrNil(key)
}

func (m *Map) Keys() []interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.keys()
}

func (m *Map) Remove(key interface{}) *Map {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.remove(key)
	return m
}

// Extend functions

func (m *Map) GetIf(key interface{}) (value interface{}, exists bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, exists = m.data[key]
	return
}

func (m *Map) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.size()
}

func (m *Map) Data() map[interface{}]interface{} {
	return m.data
}

// private

func (m *Map) set(key, value interface{}) {
	m.data[key] = value
}

func (m *Map) get(key interface{}) interface{} {
	return m.data[key]
}

func (m *Map) getOrNil(key interface{}) interface{} {
	if value, exists := m.data[key]; exists {
		return value
	}

	return nil
}

func (m *Map) size() int {
	return len(m.data)
}

func (m *Map) keys() []interface{} {
	keys := make([]interface{}, m.size())
	count := 0
	for key := range m.data {
		keys[count] = key
		count++
	}
	return keys
}

func (m *Map) remove(key interface{}) {
	delete(m.data, key)
}
