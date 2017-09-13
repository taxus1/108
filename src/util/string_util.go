package util

// Substr 截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

// SearchString 数组是否包含元素
func SearchString(list []string, x string) bool {
	if len(list) == 0 {
		return false
	}

	for _, v := range list {
		if v == x {
			return true
		}
	}
	return false
}
