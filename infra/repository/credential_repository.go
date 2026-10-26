package repository

import (
	"context"

	"github.com/DIEALOG/DIEALOG-crypto/domain"
	"github.com/DIEALOG/DIEALOG-crypto/ent"
	"github.com/DIEALOG/DIEALOG-crypto/ent/credential"
)

type credentialRepository struct {
	client *ent.Client
}

func NewCredentialRepository(client *ent.Client) domain.CredentialRepository {
	return &credentialRepository{client: client}
}

func (r *credentialRepository) Create(ctx context.Context, c *domain.Credential) error {
	exist, err := r.client.Credential.
		Query().
		Where(
			credential.AccountNameEQ(c.AccountName),
			credential.UserIDEQ(c.UserID),
		).
		Only(ctx)
	if err == nil {
		return r.client.Credential.
			UpdateOne(exist).
			SetEncryptedPassword(c.EncryptedPassword).
			SetEncryptedDek(c.EncryptedDEK).
			SetNonce(c.Nonce).
			Exec(ctx)
	}
	if !ent.IsNotFound(err) {
		return err
	}

	_, err = r.client.Credential.
		Create().
		SetUserID(c.UserID).
		SetAccountName(c.AccountName).
		SetEncryptedPassword(c.EncryptedPassword).
		SetEncryptedDek(c.EncryptedDEK).
		SetNonce(c.Nonce).
		Save(ctx)
	return err
}

func (r *credentialRepository) GetByAccountNameAndUserID(ctx context.Context, userID *domain.ID, accountName string) (*domain.Credential, error) {
	entCred, err := r.client.Credential.
		Query().
		Where(
			credential.AccountNameEQ(accountName),
			credential.UserIDEQ(*userID),
		).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return ToDomainCredential(entCred), nil
}

func (r *credentialRepository) DeleteByAccountNameAndUserID(ctx context.Context, userID *domain.ID, accountName string) error {
	n, err := r.client.Credential.
		Delete().
		Where(
			credential.AccountNameEQ(accountName),
			credential.UserIDEQ(*userID),
		).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}
