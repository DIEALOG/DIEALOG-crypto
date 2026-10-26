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
	credentials.Post("/delete", handler.DeletePassword(service))
	credentials.Delete("/delete", handler.DeletePassword(service))
	credentials.Delete("", handler.DeletePassword(service))
	credentials.Delete("/", handler.DeletePassword(service))
	credentials.Get("/:user_id/:account_name", handler.GetPassword(service))
	credentials.Delete("/:user_id/:account_name", handler.DeletePassword(service))
}
