package handler

import (
	"net/http"

	"github.com/DIEALOG/DIEALOG-crypto/internal/crypto"
	"github.com/gofiber/fiber/v3"

	"github.com/DIEALOG/DIEALOG-crypto/domain"
)

func StorePassword(credentialService domain.CredentialUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.StoreCredentialRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: "잘못된 요청 형식입니다."})
		}

		if err := credentialService.StorePassword(c.Context(), req.UserID, &req); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		return c.Status(http.StatusCreated).JSON(domain.SuccessResponse{Message: "비밀번호 저장 완료"})
	}
}

func GetPassword(credentialService domain.CredentialUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.GetPasswordRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: "잘못된 요청 형식입니다."})
		}

		// UseCase 실행 (KMS DEK 복호화 + AES-GCM 복호화)
		plainPasswordBytes, err := credentialService.GetPassword(c.Context(), req.UserID, req.AccountName)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(domain.ErrorResponse{Message: err.Error()})
		}
		defer crypto.ZeroBytes(plainPasswordBytes) // 응답 직후 메모리 파기

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"account_name": req.AccountName,
			"password":     string(plainPasswordBytes),
		})
	}
}
