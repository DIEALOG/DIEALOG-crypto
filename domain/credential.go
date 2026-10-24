package domain

import (
	"context"

	backend "github.com/DIEALOG/DIEALOG-backend/domain"
)

type Credential struct {
	ID                backend.ID
	UserID            backend.ID
	AccountName       string
	EncryptedPassword []byte
	EncryptedDEK      []byte
	Nonce             []byte
}

type StoreCredentialRequest struct {
	UserID      backend.ID `json:"user_id"`
	AccountName string     `json:"account_name"`
	Password    string     `json:"password"`
}

type GetPasswordRequest struct {
	UserID      backend.ID `json:"user_id"`
	AccountName string     `json:"account_name"`
}

type KMSService interface {
	GenerateDataKey(ctx context.Context) (plainDEK []byte, encryptedDEK []byte, err error)
	DecryptDataKey(ctx context.Context, encryptedDEK []byte) (plainDEK []byte, err error)
}

type CredentialRepository interface {
	Create(ctx context.Context, credential *Credential) error
	GetByAccountName(ctx context.Context, userID *backend.ID, accountName string) (*Credential, error)
}

type CredentialUseCase interface {
	StorePassword(ctx context.Context, userID *backend.ID, req *StoreCredentialRequest) error
	GetPassword(ctx context.Context, userID *backend.ID, accountName string) ([]byte, error)
}
