package bytes_incrementor

func Increment(bytes *[]byte, index int) {
	if index < 0 {
		temp := append([]byte{97}, *bytes...)
		*bytes = temp
	} else if (*bytes)[index] == 57 {
		oldLength := len(*bytes)
		Increment(bytes, index-1)
		if oldLength != len(*bytes) {
			index++
		}
		(*bytes)[index] = 97
		if index == 0 {
		}
	} else if (*bytes)[index] == 122 {
		(*bytes)[index] = 65
	} else if (*bytes)[index] == 90 {
		(*bytes)[index] = 48
	} else {
		(*bytes)[index] = (*bytes)[index] + 1
	}
}
