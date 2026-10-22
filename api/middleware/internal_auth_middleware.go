package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"github.com/DIEALOG/DIEALOG-crypto/config"
	"github.com/DIEALOG/DIEALOG-crypto/domain"
)

func NewInternalAuthMiddleware(env *config.Env) fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. 요청 헤더에서 시크릿 키 추출
		clientSecret := c.Get("X-Internal-Secret-Key")
		if clientSecret == "" {
			return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{
				Message: "내부 서비스 인증 헤더(X-Internal-Secret-Key)가 누락되었습니다.",
			})
		}

		// 2. 환경변수에 설정된 시크릿 키와 비교 (Constant-Time 비교로 타이밍 공격 방지)
		expectedSecret := env.InternalSecretKey
		if subtle.ConstantTimeCompare([]byte(clientSecret), []byte(expectedSecret)) != 1 {
			return c.Status(http.StatusForbidden).JSON(domain.ErrorResponse{
				Message: "유효하지 않은 내부 서비스 접근 키입니다.",
			})
		}

		return c.Next()
	}
}
