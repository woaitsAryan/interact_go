package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/validators"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App) {
	app.Post("/signup", validators.UserCreateValidator, controllers.SignUp)
	app.Post("/login", controllers.LogIn)

	userRoutes := app.Group("/users", middlewares.Protect)
	userRoutes.Get("/me", controllers.GetMe)
	userRoutes.Get("/views", controllers.GetViews)
	userRoutes.Patch("/update_password", controllers.UpdatePassword)
	userRoutes.Patch("/me", controllers.UpdateMe)
	userRoutes.Patch("/me/achievements", controllers.AddAchievement)
	userRoutes.Delete("/me/achievements/:achievementID", controllers.DeleteAchievement)
	userRoutes.Delete("/me", controllers.DeleteMe)
}
