package atcs

import (
	"atc/base"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

// Atcs  int 的包装
type Atcs map[string]interface{}

// Set 设置
func (m Atcs) Set(key string, value interface{}) interface{} {
	m[key] = value
	return nil
}

// Get 设置
func (m Atcs) Get(key string) interface{} {
	if value, ok := m[key]; ok {
		return value
	}
	return nil
}

// Add 加法
func (m Atcs) Add(key string, value int) interface{} {
	vv := m.Get(key)
	if vv == nil {
		m.Set(key, value)
	} else if v, ok := vv.(int); ok {
		m.Set(key, v+value)
	}
	return nil
}

// IsEnd 是否在结束的地方
func (m Atcs) IsEnd(index, width int) bool {
	if width == 0 {
		return false
	}
	return index != 0 && index%width == 0
}

// Others 比如10，5
func (m Atcs) Others(arr []interface{}, width int) []int {
	len := len(arr) % width
	if len != 0 {
		len = width - len
	}
	return make([]int, len, len)
}

// ArrPair 成对的数组
type ArrPair struct {
	First  []interface{}
	Second []interface{}
}

// Cut 比如10，5
func (m Atcs) Cut(arr []interface{}, width int) ArrPair {
	if len(arr) > width {
		return ArrPair{
			First:  arr[:width],
			Second: arr[width:],
		}
	}
	return ArrPair{
		First:  arr,
		Second: []interface{}{},
	}
}

// Ter 三目运算
func (m Atcs) Ter(c bool, a, b interface{}) interface{} {
	if c {
		return a
	}
	return b
}

// At 取Item
func (m Atcs) At(arr []interface{}, index int) interface{} {
	if len(arr) > index && index >= 0 {
		return arr[index]
	}
	return nil
}

// F2i float转int
func (m Atcs) F2i(value float64) int {
	return int(value)
}

// Len 数组长度
func (m Atcs) Len(arr []interface{}) int {
	return len(arr)
}

// Arr 组成数组
func (m Atcs) Arr(value ...interface{}) []interface{} {
	return value
}

// SetTo 设置到map上
func (m Atcs) SetTo(o map[string]interface{}, key string, value interface{}) interface{} {
	o[key] = value
	return nil
}

func readFileContent(path string) *string {
	f, err := os.Open(fmt.Sprintf("%s%s", base.Config.WebPath, path))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		return nil
	}
	content := string(b)
	return &content

}
func readFileContentFromCache(path string) *string {
	base.RefCacheMutex.Lock()
	defer base.RefCacheMutex.Unlock()
	if cache, ok := base.RefCache[path]; ok {
		return cache
	}
	r := readFileContent(path)
	base.RefCache[path] = r
	return r
}

// Ref 引入模块
func (m Atcs) Ref(path string) template.HTML {
	var r *string
	if base.Config.IsDev {
		r = readFileContent(path)
	} else {
		r = readFileContentFromCache(path)
	}
	if r == nil {
		return template.HTML("")
	}
	return template.HTML(*r)
}
