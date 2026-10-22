package repository

import (
	"github.com/DIEALOG/DIEALOG-crypto/domain"
	"github.com/DIEALOG/DIEALOG-crypto/ent"
)

func ToDomainCredential(e *ent.Credential) *domain.Credential {
	if e == nil {
		return nil
	}
	return &domain.Credential{
		ID:                e.ID,
		UserID:            e.UserID,
		AccountName:       e.AccountName,
		EncryptedPassword: e.EncryptedPassword,
		EncryptedDEK:      e.EncryptedDek,
		Nonce:             e.Nonce,
	}
}
