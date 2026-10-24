package main

import (
	"log"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/DIEALOG/DIEALOG-crypto/api/route"
	"github.com/DIEALOG/DIEALOG-crypto/config"
	"github.com/DIEALOG/DIEALOG-crypto/infra/database"
	"github.com/DIEALOG/DIEALOG-crypto/infra/kms"
	"github.com/DIEALOG/DIEALOG-crypto/infra/repository"
	"github.com/DIEALOG/DIEALOG-crypto/usecase"
)

func main() {
	config.NewEnv()

	app := fiber.New(fiber.Config{
		AppName:      "Fiber Ent Clean Architecture",
		ServerHeader: "Fiber",
	})

	// Use global middlewares.
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(limiter.New(limiter.Config{
		Max: 100,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(&fiber.Map{
				"status":  "fail",
				"message": "You have requested too many in a single time-frame! Please wait another minute!",
			})
		},
	}))
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	api := app.Group("/api")

	app.Get("/swagger", func(c fiber.Ctx) error {
		return c.Redirect().To("/swagger/index.html")
	})
	app.Get("/swagger/*", swaggo.HandlerDefault)
	app.Get("/docs", func(c fiber.Ctx) error {
		return c.Redirect().To("/docs/index.html")
	})
	app.Get("/docs/*", swaggo.HandlerDefault)

	client, err := database.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	kmsService := kms.NewMockKMSService([]byte(config.E.MockKMSMasterKey))
	credRepo := repository.NewCredentialRepository(client)
	credUseCase := usecase.NewCredentialUseCase(credRepo, kmsService)

	route.NewCredentialRoute(api, credUseCase)

	addr := config.E.ServerAddress
	if len(addr) > 0 && addr[0] != ':' {
		addr = ":" + addr
	}
	log.Fatal(app.Listen(addr))
}
