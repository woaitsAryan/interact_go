package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func VerificationRouter(app *fiber.App) {
	workspaceRoutes := app.Group("/verification", middlewares.Protect)
	workspaceRoutes.Get("/otp", controllers.SendVerificationCode)
	workspaceRoutes.Post("/otp", controllers.VerifyCode)
}
