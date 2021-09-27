package utils

import (
	"encoding/json"
	"sync"
)

type Array struct {
	data  []interface{}
	mutex sync.RWMutex
}

/*
Example usage
mArr := utils.NewArray()
fmt.Println(utils.Unique(mArr.Add(1,2,1,2).Values()))
*/

func NewArray() *Array {
	return new(Array)
}

func (a *Array) Add(values ...interface{}) *Array {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.add(values...)
	return a
}

/*
AddIfNotExist add one value if not existing
*/
func (a *Array) AddIfNotExist(value interface{}) bool {
	check := a.Exist(value)
	if !check {
		a.add(value)
		return true
	}
	return false
}

/*
Get a value by its index number
*/
func (a *Array) Get(index int) interface{} {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if index > a.size() {
		return nil
	}

	return a.get(index)
}

/*
RemoveNil any nil values.
Example:
	var arr = []interface{}{1, 2, 3, 4, nil, 5}
	result := RemoveNil(arr)  // [1, 2, 3, 4, 5]
*/
func (a *Array) RemoveNil(arr []interface{}) []interface{} {
	if arr == nil {
		return arr
	}
	result := make([]interface{}, 0, len(arr))
	for _, v := range arr {
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

//MinMaxInt min and max value from an array
func (a *Array) MinMaxInt() (minVal int, maxVal int) {
	max := a.data[0].(int)
	min := a.data[0].(int)
	for _, value := range a.data {
		if max < value.(int) {
			max = value.(int)
		}
		if min > value.(int) {
			min = value.(int)
		}
	}
	return min, max
}

/*
Unique a new array with duplicates removed.
Example:
	var myArray = []interface{}{1, 2, 3, 3, 4}
	result := Unique(myArray)  // [1, 2, 3, 4]
*/
func (a *Array) Unique(arr []interface{}) []interface{} {
	if arr == nil || len(arr) <= 1 {
		return arr
	}
	result := make([]interface{}, 0, len(arr))
	for _, v := range arr {
		if len(result) == 0 {
			result = append(result, v)
			continue
		}
		duplicate := false
		for _, r := range result {
			if r == v {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, v)
		}
	}
	return result
}

/*
Remove a value from the array
*/
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
