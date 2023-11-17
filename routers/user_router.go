package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/validators"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App) {
	app.Post("/signup", validators.UserCreateValidator, middlewares.EarlyAccessCheck, controllers.SignUp)
	app.Post("/login", controllers.LogIn)
	app.Post("/refresh", controllers.Refresh)

	app.Post("/early_access", controllers.GetEarlyAccessToken)

	app.Post("/recovery", controllers.SendResetURL)
	app.Post("/recovery/verify", controllers.ResetPassword)

	userRoutes := app.Group("/users", middlewares.Protect)
	userRoutes.Get("/me", controllers.GetMe)
	userRoutes.Get("/me/likes", controllers.GetMyLikes)
	userRoutes.Get("/views", controllers.GetViews)

	userRoutes.Patch("/update_password", controllers.UpdatePassword)
	userRoutes.Patch("/update_email", controllers.UpdateEmail)
	userRoutes.Patch("/update_phone_number", controllers.UpdatePhoneNo)
	userRoutes.Delete("/deactive", controllers.Deactive)

	userRoutes.Patch("/me", controllers.UpdateMe)
	userRoutes.Patch("/me/profile", controllers.EditProfile)
	userRoutes.Patch("/me/achievements", controllers.AddAchievement)
	userRoutes.Delete("/me/achievements/:achievementID", controllers.DeleteAchievement)
	userRoutes.Delete("/me", controllers.DeactivateMe)

	userRoutes.Post("/report", controllers.AddReport)
}
