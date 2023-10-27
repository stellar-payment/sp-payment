package cryptoutil

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/zenazn/pkcs7pad"
)

func AES256Encrypt(msg, iv, key []byte) (res []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to init AES chiper err: %+w", err)
	}

	encryptor := cipher.NewCBCEncrypter(block, iv)
	padded := pkcs7pad.Pad(msg, encryptor.BlockSize())

	res = make([]byte, len(padded))
	encryptor.CryptBlocks(res, padded)

	return
}

func AES256Decrypt(msg, iv, key []byte) (res []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to init AES chiper err: %+w", err)
	}

	encryptor := cipher.NewCBCDecrypter(block, iv)

	res = make([]byte, len(msg))
	encryptor.CryptBlocks(res, msg)

	res, err = pkcs7pad.Unpad(res)
	if err != nil {
		return nil, fmt.Errorf("failed to unpad plaintext err: %+w", err)
	}

	return
}

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

func VerifyHMACSHA512(msg, key, hash []byte) bool {
	digest := hmac.New(sha512.New, key)
	digest.Write(msg)
	genHash := digest.Sum(nil)

	return hmac.Equal(genHash, hash)
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

// Accept only msg and key, as IV are generated on each encryption
// rowHash are optional args to save rawbytes to be used as row-wide hash
func EncryptField(msg []byte, key []byte, rowHash *[]byte) (res []byte) {
	iv := make([]byte, 16)

	_, err := rand.Read(iv)
	if err != nil {
		// iv MUST be generated, otherwise halt op
		panic(err)
	}

	ct, err := AES256Encrypt(msg, iv, key)
	if err != nil {
		panic(err)
	}

	res = append(ct, iv...)
	if rowHash != nil {
		*rowHash = append(*rowHash, res...)
	}

	return res
}

func DecryptField(ct []byte, key []byte) (res string) {
	iv := ct[len(ct)-16:]
	ct = ct[:len(ct)-16]

	ct, err := AES256Decrypt(ct, iv, key)
	if err != nil {
		panic(err)
	}

	return string(ct)
}
