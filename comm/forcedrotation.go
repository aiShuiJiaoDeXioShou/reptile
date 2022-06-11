package comm

import "strconv"

// 把string转换为uint
func ParseUint(str string) (uint, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}