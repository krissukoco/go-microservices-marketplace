package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pwd string, cost ...int) (string, error) {
	c := 7
	for _, v := range cost {
		c = v
		break
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), c)
	return string(bytes), err
}
