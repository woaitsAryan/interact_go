package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/validators"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app *fiber.App) {
	authRoutes := app.Group("/org")
	authRoutes.Post("/signup", validators.UserCreateValidator, organization_controllers.SignUp)
	authRoutes.Post("/login", organization_controllers.LogIn)
	authRoutes.Get("/oauth/login", middlewares.ProtectRedirect, organization_controllers.OAuthLogIn)
}
