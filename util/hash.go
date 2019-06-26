package util

import (
	"crypto/sha512"
	"fmt"
)

// CalculateSha512 generates the sha512 hash from the secret
func CalculateSha512(secret string) string {
	h := sha512.New()
	h.Write([]byte(secret))
	bs := h.Sum(nil)
	hash := fmt.Sprintf("%x", bs)
	return hash
}
