package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func OauthRouter(app *fiber.App) {
	oauthRoutes := app.Group("/auth")
	oauthRoutes.Get("/signup", middlewares.ProtectRedirect, controllers.OAuthSignUp)
	oauthRoutes.Get("/login", middlewares.ProtectRedirect, controllers.OAuthLogIn)

	oauthRoutes.Get("/google", controllers.GoogleRedirect)
	oauthRoutes.Get("/google/callback", controllers.GoogleCallback)
}
