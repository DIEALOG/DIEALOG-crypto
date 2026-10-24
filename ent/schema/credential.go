package schema

import (
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Credential struct {
	ent.Schema
}

func (Credential) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().StorageKey("id"),
		field.UUID("user_id", uuid.UUID{}),
		field.String("account_name"),
		field.Bytes("encrypted_password"), // AES-256-GCM 암호문
		field.Bytes("encrypted_dek"),      // KMS KEK로 암호화된 DEK
		field.Bytes("nonce"),              // GCM Nonce (12 bytes)
	}
}
