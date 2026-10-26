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
	UserID      ID     `json:"user_id" form:"user_id" query:"user_id"`
	AccountName string `json:"account_name" form:"account_name" query:"account_name"`
	Password    string `json:"password" form:"password"`
}

type GetPasswordRequest struct {
	UserID      ID     `json:"user_id" form:"user_id" query:"user_id"`
	AccountName string `json:"account_name" form:"account_name" query:"account_name"`
}

type DeletePasswordRequest struct {
	UserID      ID     `json:"user_id" form:"user_id" query:"user_id"`
	AccountName string `json:"account_name" form:"account_name" query:"account_name"`
}

type KMSService interface {
	GenerateDataKey(ctx context.Context) (plainDEK []byte, encryptedDEK []byte, err error)
	DecryptDataKey(ctx context.Context, encryptedDEK []byte) (plainDEK []byte, err error)
}

type CredentialRepository interface {
	Create(ctx context.Context, credential *Credential) error
	GetByAccountNameAndUserID(ctx context.Context, userID *ID, accountName string) (*Credential, error)
	DeleteByAccountNameAndUserID(ctx context.Context, userID *ID, accountName string) error
}

type CredentialUseCase interface {
	StorePassword(ctx context.Context, userID *ID, req *StoreCredentialRequest) error
	GetPassword(ctx context.Context, userID *ID, accountName string) ([]byte, error)
	DeletePassword(ctx context.Context, userID *ID, accountName string) error
}
