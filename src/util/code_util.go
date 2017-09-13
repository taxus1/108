package util

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

//Krand 随机字符串
// 效率不够高
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func ToMap(v interface{}) (m map[string]interface{}, err error) {
	js, err := json.Marshal(v)
	if err != nil {
		return
	}
	if err = json.Unmarshal(js, &m); err != nil {
		return
	}
	return
}

func GenCode() string {
	now := time.Now()
	return now.Format("2006010215") + strconv.Itoa(now.Second()) + string(Krand(4, KC_RAND_KIND_NUM))
}
