package rands

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandString(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("length must be greater than 0")
	}

	letterLen := len(letterBytes)
	var strBuilder strings.Builder
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(letterLen-1)))
		if err != nil {
			return "", err
		}

		strBuilder.WriteByte(letterBytes[num.Int64()])
	}

	return strBuilder.String(), nil
}