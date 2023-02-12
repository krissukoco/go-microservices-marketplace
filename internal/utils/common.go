package utils

import (
	"math/rand"
	"time"
)

const (
	AlphaNumerics = "01234567890123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func NewAlphanumericID() string {
	count := 14
	id := ""
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		c := AlphaNumerics[rand.Intn(len(AlphaNumerics))]
		id += string(c)
	}
	return id
}
