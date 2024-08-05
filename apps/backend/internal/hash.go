package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
)

var err error

func HashString(str string) (string, error) {
	var h hash.Hash = sha256.New()
	_, err = io.WriteString(h, str)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
