package helpers

import (
	"math/rand"
	"strconv"
)

func Hash(id int) string {
	temp := strconv.Itoa(id)
	chars := "0123456789"
	code1 := make([]byte, 6)
	for i := range code1 {
		code1[i] = chars[rand.Intn(len(chars))]
	}
	code2 := make([]byte, 6)
	for i := range code2 {
		code2[i] = chars[rand.Intn(len(chars))]
	}

	return string(code1) + temp + string(code2)
}

func Unhash(id string) string {
	return id[6 : len(id)-6]
}
