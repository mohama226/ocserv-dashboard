package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func CalculateSHA256(r io.Reader) (string, error) {

	hash := sha256.New()

	if _, err := io.Copy(hash, r); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
