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
			credential.AccountIDEQ(c.AccountID),
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
		SetAccountID(c.AccountID).
		SetEncryptedPassword(c.EncryptedPassword).
		SetEncryptedDek(c.EncryptedDEK).
		SetNonce(c.Nonce).
		Save(ctx)
	return err
}

func (r *credentialRepository) GetByAccountIDAndUserID(ctx context.Context, userID *domain.ID, accountID *domain.ID) (*domain.Credential, error) {
	entCred, err := r.client.Credential.
		Query().
		Where(
			credential.AccountIDEQ(*accountID),
			credential.UserIDEQ(*userID),
		).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return ToDomainCredential(entCred), nil
}

func (r *credentialRepository) DeleteByAccountIDAndUserID(ctx context.Context, userID *domain.ID, accountID *domain.ID) error {
	n, err := r.client.Credential.
		Delete().
		Where(
			credential.AccountIDEQ(*accountID),
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
