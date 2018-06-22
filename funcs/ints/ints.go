package ints

import (
	"log"
	"strconv"
)

// Ints  int 的包装
type Ints struct{}

// Add 加法
func (m Ints) Add(a int, b int) int {
	return a + b
}

// Sub 减法
func (m Ints) Sub(a int, b int) int {
	return a - b
}

// Mod 取模
func (m Ints) Mod(a int, b int) int {
	return a % b
}

// Arr 创建数组
func (m Ints) Arr(len int) []int {
	return make([]int, len, len)
}

// Int to int
func (m Ints) Int(v interface{}) int {
	switch i := v.(type) {
	case int64:
		return int(i)
	case int:
		return i
	case float64:
		return int(i)
	case string:
		r, err := strconv.Atoi(i)
		if err != nil {
			log.Println(err)
		}
		return r
	case int32:
		return int(i)
	case int16:
		return int(i)
	case int8:
		return int(i)
	case float32:
		return int(i)
	case bool:
		if i {
			return 1
		}
	case uint:
		return int(i)
	case uint64:
		return int(i)
	case uint32:
		return int(i)
	case uint16:
		return int(i)
	case uint8:
		return int(i)
	}
	return 0
}

// Int64 to int64
func (m Ints) Int64(v interface{}) int64 {
	switch i := v.(type) {
	case int:
		return int64(i)
	case int64:
		return i
	case float64:
		return int64(i)
	case string:
		r, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			log.Println(err)
		}
		return r
	case int32:
		return int64(i)
	case int16:
		return int64(i)
	case int8:
		return int64(i)
	case float32:
		return int64(i)
	case bool:
		if i {
			return 1
		}
	case uint:
		return int64(i)
	case uint64:
		return int64(i)
	case uint32:
		return int64(i)
	case uint16:
		return int64(i)
	case uint8:
		return int64(i)
	}
	return 0
}

// At 取Item
func (m Ints) At(arr []int, index int) int {
	if len(arr) > index && index >= 0 {
		return arr[index]
	}
	return 0
}
