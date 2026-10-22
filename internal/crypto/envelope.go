package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func EncryptAESGCM(plainDEK, plaintext, aad []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(plainDEK)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, aad)
	return ciphertext, nonce, nil
}

func DecryptAESGCM(plainDEK, ciphertext, nonce, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(plainDEK)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, errors.New("데이터 변조 검증 실패: 암호문이 유효하지 않거나 UserID가 일치하지 않습니다")
	}

	return plaintext, nil
}
