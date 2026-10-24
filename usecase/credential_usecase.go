package usecase

import (
	"context"

	backend "github.com/DIEALOG/DIEALOG-backend/domain"
	"github.com/DIEALOG/DIEALOG-crypto/domain"
	"github.com/DIEALOG/DIEALOG-crypto/ent"
	"github.com/DIEALOG/DIEALOG-crypto/internal/crypto"
)

type credentialUseCase struct {
	repo domain.CredentialRepository
	kms  domain.KMSService
}

func NewCredentialUseCase(repo domain.CredentialRepository, kms domain.KMSService) domain.CredentialUseCase {
	return &credentialUseCase{
		repo: repo,
		kms:  kms,
	}
}

func (u *credentialUseCase) StorePassword(ctx context.Context, userID *backend.ID, req *domain.StoreCredentialRequest) error {
	// 1. 메모리 소거용 byte 슬라이스
	plainPassword := []byte(req.Password)
	defer crypto.ZeroBytes(plainPassword)

	// 2. KMS에서 DEK 생성
	plainDEK, encryptedDEK, err := u.kms.GenerateDataKey(ctx)
	if err != nil {
		return err
	}
	defer crypto.ZeroBytes(plainDEK) // 사용 후 즉시 파기

	// 3. UserID를 AAD(Authenticated Associated Data)로 함께 바인딩하여 암호화
	aad := (*userID)[:]
	ciphertext, nonce, err := crypto.EncryptAESGCM(plainDEK, plainPassword, aad)
	if err != nil {
		return err
	}

	// 4. DB 저장
	cred := &domain.Credential{
		UserID:            *userID,
		AccountName:       req.AccountName,
		EncryptedPassword: ciphertext,
		EncryptedDEK:      encryptedDEK,
		Nonce:             nonce,
	}

	return u.repo.Create(ctx, cred)
}

func (u *credentialUseCase) GetPassword(ctx context.Context, userID *backend.ID, accountName string) ([]byte, error) {
	cred, err := u.repo.GetByAccountName(ctx, userID, accountName)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewNotFoundError(err)
		}
		return nil, err
	}

	// KMS로 DEK 복호화
	plainDEK, err := u.kms.DecryptDataKey(ctx, cred.EncryptedDEK)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}
	defer crypto.ZeroBytes(plainDEK)

	// AES-GCM 복호화 (UserID AAD 검증)
	aad := (*userID)[:]
	pwd, err := crypto.DecryptAESGCM(plainDEK, cred.EncryptedPassword, cred.Nonce, aad)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return pwd, nil
}
