package domain

import (
	"context"
)

type Credential struct {
	ID                ID
	UserID            ID
	AccountName       string
	EncryptedPassword []byte
	EncryptedDEK      []byte
	Nonce             []byte
}

type StoreCredentialRequest struct {
	UserID      ID     `json:"user_id"`
	AccountName string `json:"account_name"`
	Password    string `json:"password"`
}

type GetPasswordRequest struct {
	UserID      ID     `json:"user_id"`
	AccountName string `json:"account_name"`
}

type KMSService interface {
	GenerateDataKey(ctx context.Context) (plainDEK []byte, encryptedDEK []byte, err error)
	DecryptDataKey(ctx context.Context, encryptedDEK []byte) (plainDEK []byte, err error)
}

type CredentialRepository interface {
	Create(ctx context.Context, credential *Credential) error
	GetByAccountName(ctx context.Context, userID *ID, accountName string) (*Credential, error)
}

type CredentialUseCase interface {
	StorePassword(ctx context.Context, userID *ID, req *StoreCredentialRequest) error
	GetPassword(ctx context.Context, userID *ID, accountName string) ([]byte, error)
}
