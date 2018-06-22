package floats

import "strconv"

// Floats  float 的包装
type Floats struct{}

// Add 加法
func (m Floats) Add(a float64, b float64) float64 {
	return a + b
}

// Sub 减法
func (m Floats) Sub(a float64, b float64) float64 {
	return a - b
}

// Fixed toString
func (m Floats) Fixed(a float64, l int) string {
	return strconv.FormatFloat(a, 'f', l, 64)
}

// At 取Item
func (m Floats) At(arr []float64, index int) float64 {
	if len(arr) > index && index >= 0 {
		return arr[index]
	}
	return 0.0
}
