package route

import (
	"github.com/gofiber/fiber/v3"

	"github.com/DIEALOG/DIEALOG-crypto/api/handler"
	"github.com/DIEALOG/DIEALOG-crypto/api/middleware"
	"github.com/DIEALOG/DIEALOG-crypto/config"
	"github.com/DIEALOG/DIEALOG-crypto/domain"
)

func NewCredentialRoute(app fiber.Router, service domain.CredentialUseCase) {
	credentials := app.Group("/credentials")
	credentials.Use(middleware.NewInternalAuthMiddleware(&config.E))

	credentials.Post("/store", handler.StorePassword(service))
	credentials.Post("/retrieve", handler.GetPassword(service))
}
