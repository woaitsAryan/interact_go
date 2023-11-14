package main

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.AddLogger()
	initializers.ConnectToCache()
	initializers.AutoMigrate()

	config.InitializeOAuthGoogle()
}

func main() {
	defer initializers.LoggerCleanUp()
	app := fiber.New(fiber.Config{
		ErrorHandler: helpers.ErrorHandler,
		BodyLimit:    config.BODY_LIMIT,
	})

	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(config.CORS())
	app.Use(config.RATE_LIMITER())
	// app.Use(config.API_CHECKER)

	app.Static("/", "./public")

	routers.Config(app)

	app.Listen(":" + initializers.CONFIG.PORT)

}
