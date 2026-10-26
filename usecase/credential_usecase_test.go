package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DIEALOG/DIEALOG-crypto/domain"
	"github.com/DIEALOG/DIEALOG-crypto/ent"
	"github.com/DIEALOG/DIEALOG-crypto/usecase"
	"github.com/google/uuid"
)

type mockRepo struct {
	createFunc                      func(ctx context.Context, credential *domain.Credential) error
	getByAccountNameAndUserIDFunc   func(ctx context.Context, userID *domain.ID, accountName string) (*domain.Credential, error)
	deleteByAccountNameAndUserIDFunc func(ctx context.Context, userID *domain.ID, accountName string) error
}

func (m *mockRepo) Create(ctx context.Context, credential *domain.Credential) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, credential)
	}
	return nil
}

func (m *mockRepo) GetByAccountNameAndUserID(ctx context.Context, userID *domain.ID, accountName string) (*domain.Credential, error) {
	if m.getByAccountNameAndUserIDFunc != nil {
		return m.getByAccountNameAndUserIDFunc(ctx, userID, accountName)
	}
	return nil, nil
}

func (m *mockRepo) DeleteByAccountNameAndUserID(ctx context.Context, userID *domain.ID, accountName string) error {
	if m.deleteByAccountNameAndUserIDFunc != nil {
		return m.deleteByAccountNameAndUserIDFunc(ctx, userID, accountName)
	}
	return nil
}

type mockKMS struct {
	generateDataKeyFunc func(ctx context.Context) ([]byte, []byte, error)
	decryptDataKeyFunc  func(ctx context.Context, encryptedDEK []byte) ([]byte, error)
}

func (m *mockKMS) GenerateDataKey(ctx context.Context) ([]byte, []byte, error) {
	if m.generateDataKeyFunc != nil {
		return m.generateDataKeyFunc(ctx)
	}
	return []byte("12345678901234567890123456789012"), []byte("encDEK"), nil
}

func (m *mockKMS) DecryptDataKey(ctx context.Context, encryptedDEK []byte) ([]byte, error) {
	if m.decryptDataKeyFunc != nil {
		return m.decryptDataKeyFunc(ctx, encryptedDEK)
	}
	return []byte("12345678901234567890123456789012"), nil
}

func TestDeletePassword_Success(t *testing.T) {
	userID := uuid.New()
	accountName := "google"

	repo := &mockRepo{
		deleteByAccountNameAndUserIDFunc: func(ctx context.Context, uid *domain.ID, acc string) error {
			if *uid != userID || acc != accountName {
				t.Errorf("unexpected arguments: uid=%v, acc=%v", uid, acc)
			}
			return nil
		},
	}
	kms := &mockKMS{}
	uc := usecase.NewCredentialUseCase(repo, kms)

	err := uc.DeletePassword(context.Background(), &userID, accountName)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDeletePassword_NotFound(t *testing.T) {
	userID := uuid.New()
	accountName := "google"

	repo := &mockRepo{
		deleteByAccountNameAndUserIDFunc: func(ctx context.Context, uid *domain.ID, acc string) error {
			return &ent.NotFoundError{}
		},
	}
	kms := &mockKMS{}
	uc := usecase.NewCredentialUseCase(repo, kms)

	err := uc.DeletePassword(context.Background(), &userID, accountName)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var domainErr *domain.Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain.Error, got %T", err)
	}
	if domainErr.StatusCode != 404 {
		t.Fatalf("expected status 404, got %d", domainErr.StatusCode)
	}
}
