package domain

import (
	"context"
)

type Credential struct {
	ID                ID
	UserID            ID
	AccountID         ID
	EncryptedPassword []byte
	EncryptedDEK      []byte
	Nonce             []byte
}

type StoreCredentialRequest struct {
	UserID    ID     `json:"user_id" form:"user_id" query:"user_id"`
	AccountID ID     `json:"account_id" form:"account_id" query:"account_id"`
	Password  string `json:"password" form:"password"`
}

type GetPasswordRequest struct {
	UserID    ID `json:"user_id" form:"user_id" query:"user_id"`
	AccountID ID `json:"account_id" form:"account_id" query:"account_id"`
}

type DeletePasswordRequest struct {
	UserID    ID `json:"user_id" form:"user_id" query:"user_id"`
	AccountID ID `json:"account_id" form:"account_id" query:"account_id"`
}

type KMSService interface {
	GenerateDataKey(ctx context.Context) (plainDEK []byte, encryptedDEK []byte, err error)
	DecryptDataKey(ctx context.Context, encryptedDEK []byte) (plainDEK []byte, err error)
}

type CredentialRepository interface {
	Create(ctx context.Context, credential *Credential) error
	GetByAccountIDAndUserID(ctx context.Context, userID *ID, accountID *ID) (*Credential, error)
	DeleteByAccountIDAndUserID(ctx context.Context, userID *ID, accountID *ID) error
}

type CredentialUseCase interface {
	StorePassword(ctx context.Context, userID *ID, req *StoreCredentialRequest) error
	GetPassword(ctx context.Context, userID *ID, accountID *ID) ([]byte, error)
	DeletePassword(ctx context.Context, userID *ID, accountID *ID) error
}
