package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/validators"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app *fiber.App) {
	oauthRoutes := app.Group("/org")
	oauthRoutes.Post("/signup", validators.UserCreateValidator, organization_controllers.SignUp)
	app.Post("/login", organization_controllers.LogIn)
}
