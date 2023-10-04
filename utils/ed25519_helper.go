package utils

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"os"
	"strings"
)

/*
to generate PEM key:

openssl genpkey -algorithm ed25519 -outform PEM -out private.pem

to extract public key:

openssl pkey -in private.pem -pubout > public.pem

*/

// DecodePublicKey decodes base64-coded public key into ed25519 key
func DecodePublicKey(public string) (*[32]byte, error) {

	var pubKey [32]byte
	pKey, err := base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, fmt.Errorf("error decoding public string %s ", public)
	}
	rkey, err := x509.ParsePKIXPublicKey(pKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding public string %s ", public)
	}
	copy(pubKey[:], rkey.(ed25519.PublicKey))
	return &pubKey, nil
}

// DecodePrivateKey decodes base64-coded private key into ed25519 key
func DecodePrivateKey(private string) (*[64]byte, error) {
	var prvKey [64]byte

	pKey, err := base64.StdEncoding.DecodeString(private)
	if err != nil {
		return nil, fmt.Errorf("error decoding private key")
	}

	realpKkey, err := x509.ParsePKCS8PrivateKey(pKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding private key")
	}
	copy(prvKey[:], realpKkey.(ed25519.PrivateKey))
	return &prvKey, nil

}

// ReadPemPrivateKey reads ed25519 private key from the file, PEM format, and decodes it
func ReadPemPrivateKey(filename string) (*[64]byte, error) {

	prefix := "BEGIN PRIVATE KEY-----"
	suffix := "-----END PRIVATE KEY"
	var pKey []byte
	var err error

	pKey, err = os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s : %v", filename, err.Error())
	}

	pkS := strings.Replace(string(pKey), "\n", "", -1)
	ndx := strings.Index(pkS, prefix)
	if ndx < 0 {
		return nil, fmt.Errorf("error reading file %s", filename)
	}
	pkS = pkS[ndx+len(prefix):]
	ndx = strings.Index(pkS, suffix)
	if ndx < 0 {
		return nil, fmt.Errorf("error reading file %s", filename)
	}
	pkS = pkS[:ndx]
	return DecodePrivateKey(pkS)
}
