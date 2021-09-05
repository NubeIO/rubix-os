package utils

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Array struct {
	data  []interface{}
	mutex sync.RWMutex
}

func NewArray() *Array {
	return new(Array)
}

func (a *Array) Add(values ...interface{}) *Array {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.add(values...)

	return a
}

//AddIfNotExist add one value if not existing
func (a *Array) AddIfNotExist(value interface{})  bool {
	check := a.Exist(value)
	if !check {
		fmt.Println(check)
		a.add(value)
		return true
	}
	return false
}


func (a *Array) Get(index int) interface{} {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if index > a.size() {
		return nil
	}

	return a.get(index)
}

func (a *Array) Remove(value interface{}) *Array {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	i := a.index(value)

	if i == -1 {
		return a
	}

	a.data = a.data[:i+copy((a.data)[i:], (a.data)[i+1:])]

	return a
}

func (a *Array) Index(value interface{}) int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.index(value)
}

func (a *Array) Exist(value interface{}) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.index(value) > -1
}

func (a Array) Includes(values ...interface{}) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	var count = 0

	for _, value := range values {
		if a.index(value) > -1 {
			count++
		}
	}

	return len(values) == count
}

func (a *Array) Values() []interface{} {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	newArr := make([]interface{}, len(a.data), len(a.data))
	copy(newArr, a.data[:])
	return newArr
}

// Data unsafe function, return data
func (a *Array) Data() []interface{} {
	return a.data
}

func (a *Array) Size() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.size()
}

func (a Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.data)
}

func (a *Array) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &a.data)
}



// private

func (a *Array) get(index int) interface{} {
	return a.data[index]
}

func (a *Array) add(values ...interface{}) {
	a.data = append(a.data, values...)
}

func (a *Array) size() int {

	return len(a.data)
}

func (a *Array) index(value interface{}) int {
	if a.size() == 0 {
		return -1
	}

	for _index, _value := range a.data {
		if _value == value {
			return _index
		}
	}

	return -1
}

