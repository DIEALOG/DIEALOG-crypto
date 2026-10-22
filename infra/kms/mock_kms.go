package kms

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/DIEALOG/DIEALOG-crypto/domain"
)

type mockKMSService struct {
	masterKey []byte // KEK 역할을 하는 32bytes 로컬 마스터 키
}

func NewMockKMSService(masterKey []byte) domain.KMSService {
	if len(masterKey) != 32 {
		panic("Mock KMS 마스터 키는 정확히 32바이트(AES-256)여야 합니다.")
	}
	return &mockKMSService{masterKey: masterKey}
}

// GenerateDataKey 1. 랜덤 32byte 평문 DEK 생성 후, 로컬 마스터 키로 암호화하여 (plainDEK, encryptedDEK) 반환
func (m *mockKMSService) GenerateDataKey(ctx context.Context) ([]byte, []byte, error) {
	// 1-1. 무작위 Plain DEK 생성 (32 bytes)
	plainDEK := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, plainDEK); err != nil {
		return nil, nil, fmt.Errorf("DEK 생성 실패: %w", err)
	}

	// 1-2. 마스터 키(KEK)로 Plain DEK 암호화
	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	// Nonce + Ciphertext를 합쳐서 encryptedDEK로 리턴
	encryptedDEK := gcm.Seal(nonce, nonce, plainDEK, nil)

	return plainDEK, encryptedDEK, nil
}

// DecryptDataKey 2. encryptedDEK를 로컬 마스터 키로 복호화하여 plainDEK 반환
func (m *mockKMSService) DecryptDataKey(ctx context.Context, encryptedDEK []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedDEK) < nonceSize {
		return nil, errors.New("유효하지 않은 encryptedDEK 길이입니다")
	}

	nonce, ciphertext := encryptedDEK[:nonceSize], encryptedDEK[nonceSize:]
	plainDEK, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("DEK 복호화 실패: %w", err)
	}

	return plainDEK, nil
}
