package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash"
)

// GenerateHMACBase64 generates base64 encoded signature of  HMAC  256, 384, or 512
func GenerateHMACBase64(algo int16, secretKey string, d string) string {
	hmac := generateHMAC(algo, secretKey, d)
	return base64.StdEncoding.EncodeToString(hmac)
}

// GenerateHMACHexEncoded generates hex encoded signature of  HMAC  256, 384, or 512
func GenerateHMACHexEncoded(algo int16, secretKey string, d string) string {
	hmac := generateHMAC(algo, secretKey, d)
	return hex.EncodeToString(hmac)
}

// generateHMAC generates HMAC signature.
func generateHMAC(algo int16, secretKey string, d string) []byte {
	var h hash.Hash
	if algo == 256 {
		h = hmac.New(sha256.New, []byte(secretKey))
	} else if algo == 384 {
		h = hmac.New(sha512.New384, []byte(secretKey))
	} else if algo == 512 {
		h = hmac.New(sha512.New, []byte(secretKey))
	}
	h.Write([]byte(d))
	return h.Sum(nil)
}
