package utils

import (
	"crypto/rand"
	"math/big"
)

const alph = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomCode(length int) (string, error) {
	res := make([]byte, length)
	n := int64(len(alph))
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(n))
		if err != nil {
			return "", err
		}
		res[i] = alph[num.Int64()]
	}
	return string(res), nil
}