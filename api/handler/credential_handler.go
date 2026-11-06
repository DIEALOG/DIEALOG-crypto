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

		if req.UserID == (domain.ID{}) || req.AccountID == (domain.ID{}) {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: "잘못된 요청 형식입니다."})
		}

		if err := credentialService.StorePassword(c.RequestCtx(), &req.UserID, &req); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		return c.Status(http.StatusCreated).JSON(domain.SuccessResponse{Message: "비밀번호 저장 완료"})
	}
}

func GetPassword(credentialService domain.CredentialUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.GetPasswordRequest
		if len(c.Body()) > 0 {
			if err := c.Bind().Body(&req); err != nil {
				return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: "잘못된 요청 형식입니다."})
			}
		}
		if req.UserID == (domain.ID{}) || req.AccountID == (domain.ID{}) {
			_ = c.Bind().Query(&req)
		}
		if req.UserID == (domain.ID{}) || req.AccountID == (domain.ID{}) {
			if uidStr := c.Params("user_id"); uidStr != "" {
				if parsedUID, err := domain.StringToID(uidStr); err == nil {
					req.UserID = parsedUID
				}
			}
			if req.AccountID == (domain.ID{}) {
				if aidStr := c.Params("account_id"); aidStr != "" {
					if parsedAID, err := domain.StringToID(aidStr); err == nil {
						req.AccountID = parsedAID
					}
				}
			}
		}

		if req.UserID == (domain.ID{}) || req.AccountID == (domain.ID{}) {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: "잘못된 요청 형식입니다."})
		}

		// UseCase 실행 (KMS DEK 복호화 + AES-GCM 복호화)
		plainPasswordBytes, err := credentialService.GetPassword(c.RequestCtx(), &req.UserID, &req.AccountID)
		if err != nil {
			return RespondError(c, err)
		}
		defer crypto.ZeroBytes(plainPasswordBytes) // 응답 직후 메모리 파기

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"account_id": req.AccountID,
			"password":   string(plainPasswordBytes),
		})
	}
}

func DeletePassword(credentialService domain.CredentialUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.DeletePasswordRequest
		if len(c.Body()) > 0 {
			if err := c.Bind().Body(&req); err != nil {
				return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: "잘못된 요청 형식입니다."})
			}
		}
		if req.UserID == (domain.ID{}) || req.AccountID == (domain.ID{}) {
			_ = c.Bind().Query(&req)
		}
		if req.UserID == (domain.ID{}) || req.AccountID == (domain.ID{}) {
			if uidStr := c.Params("user_id"); uidStr != "" {
				if parsedUID, err := domain.StringToID(uidStr); err == nil {
					req.UserID = parsedUID
				}
			}
			if req.AccountID == (domain.ID{}) {
				if aidStr := c.Params("account_id"); aidStr != "" {
					if parsedAID, err := domain.StringToID(aidStr); err == nil {
						req.AccountID = parsedAID
					}
				}
			}
		}

		if req.UserID == (domain.ID{}) || req.AccountID == (domain.ID{}) {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: "잘못된 요청 형식입니다."})
		}

		if err := credentialService.DeletePassword(c.RequestCtx(), &req.UserID, &req.AccountID); err != nil {
			return RespondError(c, err)
		}

		return c.Status(http.StatusOK).JSON(domain.SuccessResponse{Message: "비밀번호 삭제 완료"})
	}
}
