package rands

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"strings"
)

type Randoms interface {
	String(n int, charSet string) (string, error)
}

type randoms struct {
	source io.Reader
}

const (
	DefaultCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

var (
	DefaultSource = rand.Reader
)

func New(source io.Reader) Randoms {
	if source == nil {
		source = rand.Reader
	}

	return randoms{
		source: source,
	}
}

func (rs randoms) String(n int, charSet string) (string, error) {
	if n <= 0 {
		return "", errors.New("length must be greater than 0")
	}

	letterLen := len(charSet)
	var strBuilder strings.Builder
	for i := 0; i < n; i++ {
		num, err := rand.Int(rs.source, big.NewInt(int64(letterLen-1)))
		if err != nil {
			return "", err
		}

		strBuilder.WriteByte(charSet[num.Int64()])
	}

	return strBuilder.String(), nil
}
