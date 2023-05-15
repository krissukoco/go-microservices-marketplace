package utils

import (
	"math/rand"
	"time"
)

const (
	AlphaNumerics = "01234567890123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func NewAlphanumericID(count ...int) string {
	c := 14
	for _, v := range count {
		c = v
		break
	}
	id := ""
	for i := 0; i < c; i++ {
		rand.Seed(time.Now().UnixNano())
		c := AlphaNumerics[rand.Intn(len(AlphaNumerics))]
		id += string(c)
	}
	return id
}
