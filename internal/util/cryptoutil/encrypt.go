package cryptoutil

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// Implementation of Java's SHA256WithRSA in Go
// Algorithm used are RSA-PKCS#1v1.5 and SHA256
func SignSHA256WithRSA(msg []byte, key *rsa.PrivateKey) (res []byte, err error) {
	digest := sha256.Sum256(msg)

	res, err = rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt hash err: %+v", err)
	}

	return
}

func VerifySHA256WithRSA(msg []byte, key *rsa.PublicKey, signature []byte) (err error) {
	digest := sha256.New()
	digest.Write(msg)

	err = rsa.VerifyPKCS1v15(key, crypto.SHA256, digest.Sum(nil), signature)
	if err != nil {
		return fmt.Errorf("failed to encrypt hash err: %+v", err)
	}

	return
}

func HMACSHA512(msg []byte, key []byte) []byte {
	digest := hmac.New(sha512.New, key)
	digest.Write(msg)
	return digest.Sum(nil)
}

func LoadPublicKey(key []byte) (pk *rsa.PublicKey, err error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("failed to decode public key err: no public key found")
	}

	keyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key err: %+v", err)
	}

	pk = keyInterface.(*rsa.PublicKey)

	return
}
