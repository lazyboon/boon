package cid

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func encrypt(data []byte, key []byte) (string, error) {
	// align key
	var err error
	key, err = alignKey(key)
	if err != nil {
		return "", err
	}

	// new cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	size := block.BlockSize()

	// pkcs7 padding
	padding := size - len(data)%size
	data = append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// aes cbc encrypt
	mode := cipher.NewCBCEncrypter(block, key[:size])
	encrypted := make([]byte, len(data))
	mode.CryptBlocks(encrypted, data)

	// base64 encode
	return base64.RawURLEncoding.EncodeToString(encrypted), nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
	// align key
	var err error
	key, err = alignKey(key)
	if err != nil {
		return nil, err
	}

	// base64 decode
	data, err = base64.RawURLEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	// aes cbc decrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	orig := make([]byte, len(data))
	mode.CryptBlocks(orig, data)

	// un padding
	length := len(orig)
	return orig[:(length - int(orig[length-1]))], nil
}

func alignKey(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("the key cannot be empty")
	}
	l := len(key)
	var ans []byte
	if l <= 16 {
		ans = make([]byte, 16)
	} else if l <= 24 {
		ans = make([]byte, 24)
	} else if l <= 32 {
		ans = make([]byte, 32)
	} else {
		return nil, errors.New("the key cannot exceed 32 bytes")
	}
	for i, d := range key {
		if i > 31 {
			break
		}
		ans[i] = d
	}
	return ans, nil
}
