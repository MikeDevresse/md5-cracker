package service

import "strings"

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Convert10to62 Convert an integer number to base62 as string
func Convert10to62(number int) string {
	res := ""
	for {
		res = string(ALPHABET[number%62]) + res
		number = number / 62
		if number == 0 {
			break
		}
	}
	return res
}

// Convert62to10 Convert a base 62 string to base10 integer
func Convert62to10(number string) int {
	sum := strings.Index(ALPHABET, string(number[0]))
	base := len(ALPHABET)
	for i := 0; i < len(number); i++ {
		sum = base*sum + strings.Index(ALPHABET, string(number[i]))
	}
	return sum
}
