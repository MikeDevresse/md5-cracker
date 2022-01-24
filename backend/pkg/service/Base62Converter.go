package service

import "strings"

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Convert10to62(number int) string {
	res := ""
	for {
		res = string(alphabet[number%62]) + res
		number = number / 62
		if number == 0 {
			break
		}
	}
	return res
}

func Convert62to10(number string) int {
	sum := strings.Index(alphabet, string(number[0]))
	base := len(alphabet)
	for i := 0; i < len(number); i++ {
		sum = base*sum + strings.Index(alphabet, string(number[i]))
	}
	return sum
}
