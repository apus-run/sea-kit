package slice

import "github.com/apus-run/sea-kit/slice/internal/slice"

// Delete 删除 index 处的元素
func Delete[Src any](src []Src, index int) ([]Src, error) {
	res, _, err := slice.Delete[Src](src, index)
	return res, err
}

// FilterDelete 删除符合条件的元素
// 考虑到性能问题，所有操作都会在原切片上进行
// 被删除元素之后的元素会往前移动，有且只会移动一次
func FilterDelete[Src any](src []Src, m func(idx int, src Src) bool) []Src {
	// 记录被删除的元素位置，也称空缺的位置
	emptyPos := 0
	for idx := range src {
		// 判断是否满足删除的条件
		if m(idx, src[idx]) {
			continue
		}
		// 移动元素
		src[emptyPos] = src[idx]
		emptyPos++
	}
	return src[:emptyPos]
}
