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
	_, err := r.client.Credential.
		Create().
		SetUserID(c.UserID).
		SetAccountName(c.AccountName).
		SetEncryptedPassword(c.EncryptedPassword).
		SetEncryptedDek(c.EncryptedDEK).
		SetNonce(c.Nonce).
		Save(ctx)
	return err
}

func (r *credentialRepository) GetByAccountName(ctx context.Context, userID *domain.ID, accountName string) (*domain.Credential, error) {
	entCred, err := r.client.Credential.
		Query().
		Where(
			credential.UserIDEQ(*userID),
			credential.AccountNameEQ(accountName),
		).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return ToDomainCredential(entCred), nil
}
