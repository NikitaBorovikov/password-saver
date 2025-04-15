package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func Encrypt(plainText, key []byte) (string, error) {
	// Creating an AES cipher with a key 32 bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generating a random IV with a length of 12 bytes
	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Returns the block cipher, converted to GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := aesgcm.Seal(nil, iv, plainText, nil)

	// Encoding the result in Base64 (IV + ciphertext)
	encText := base64.StdEncoding.EncodeToString(append(iv, ciphertext...))

	return encText, nil
}

func Decrypt(ciphertext []byte, key []byte) (string, error) {
	// Creating an AES cipher with a key 32 bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Checking the length of the data (there must be at least a IV)
	if len(ciphertext) < 12 {
		return "", errors.New("ciphertext too short")
	}
	// Extracting the IV
	iv := ciphertext[:12]
	ciphertext = ciphertext[12:]

	// Returns the block cipher, converted to GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Decrypting and verifying authentication
	plainText, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
