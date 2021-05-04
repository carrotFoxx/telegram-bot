package host_manager

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"fmt"
)

func Encrypt(key []byte, plaintext string) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Errorf("NewCipher(%d bytes) = %s", len(key), err)
	}

	dd := len(plaintext) % 16
	var bytePlaintext []byte
	if dd != 0 {
		bytePlaintext = []byte(plaintext)

		for i := 0; i < 16-dd; i++ {
			bytePlaintext = append(bytePlaintext, 0x00)
		}

	}

	var n int = len(bytePlaintext) / 16
	var out []byte
	for i := 0; i < n; i++ {
		temp := make([]byte, 16)
		c.Encrypt(temp, bytePlaintext[i*16:(i+1)*16])
		out = append(out, temp...)
	}
	return out
}

func DecryptFromString(key []byte, ct string) []byte {
	ciphertext, _ := hex.DecodeString(ct)
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Errorf("NewCipher(%d bytes) = %s", len(key), err)
		//panic(err)
	}
	plain := make([]byte, len(ciphertext))
	n := len(ciphertext) / 16
	for i := 0; i < n; i++ {
		temp := make([]byte, 16)
		c.Decrypt(temp, ciphertext[i*16:(i+1)*16])
		plain = append(plain, temp...)
	}
	var b bytes.Buffer
	for _, r := range string(plain[:]) {
		if r != 0x00 {
			b.WriteRune(r)
		}
	}
	return b.Bytes()
}

func Decrypt(key []byte, ciphertext []byte) []byte {

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Errorf("NewCipher(%d bytes) = %s", len(key), err)
		//panic(err)
	}
	plain := make([]byte, len(ciphertext))
	n := len(ciphertext) / 16
	for i := 0; i < n; i++ {
		temp := make([]byte, 16)
		c.Decrypt(temp, ciphertext[i*16:(i+1)*16])
		plain = append(plain, temp...)
	}
	var b bytes.Buffer
	for _, r := range string(plain[:]) {
		if r != 0x00 {
			b.WriteRune(r)
		}
	}
	return b.Bytes()
}
