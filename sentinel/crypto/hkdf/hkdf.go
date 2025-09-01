package hkdf

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"hash"
)

// HKDF implements HMAC-based Key Derivation Function as specified in RFC 5869
type HKDF struct {
	hashFunc func() hash.Hash
	secret   []byte
	salt     []byte
	info     []byte
}

// New creates a new HKDF instance
func New(secret, salt, info []byte) *HKDF {
	return &HKDF{
		hashFunc: sha512.New,
		secret:   secret,
		salt:     salt,
		info:     info,
	}
}

// DeriveKey derives a key of the specified length using HKDF
func (h *HKDF) DeriveKey(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("key length must be positive")
	}

	// Extract step
	prk := h.extract()

	// Expand step
	key := h.expand(prk, length)

	return key, nil
}

// extract performs the HKDF extract step
func (h *HKDF) extract() []byte {
	if h.salt == nil {
		h.salt = make([]byte, h.hashFunc().Size())
	}
	mac := hmac.New(h.hashFunc, h.salt)
	mac.Write(h.secret)
	return mac.Sum(nil)
}

// expand performs the HKDF expand step
func (h *HKDF) expand(prk []byte, length int) []byte {
	hashLen := h.hashFunc().Size()
	n := (length + hashLen - 1) / hashLen // ceil(length / hashLen)

	var okm []byte
	var t []byte

	for i := 1; i <= n; i++ {
		mac := hmac.New(h.hashFunc, prk)
		mac.Write(t)
		mac.Write(h.info)
		mac.Write([]byte{byte(i)})
		t = mac.Sum(nil)
		okm = append(okm, t...)
	}

	return okm[:length]
}

// DeriveKey is a convenience function to derive a key using HKDF
func DeriveKey(secret, salt, info []byte, length int) ([]byte, error) {
	hkdf := New(secret, salt, info)
	return hkdf.DeriveKey(length)
}
