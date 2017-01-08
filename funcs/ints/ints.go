package ints

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
	return a - b
}

// Arr 创建数组
func (m Ints) Arr(len int) []int {
	return make([]int, len, len)
}

// Int 创建数组
func (m Ints) Int(v interface{}) int {
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}
