package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/auth_controllers"
	"github.com/Pratham-Mishra04/interact/controllers/user_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/validators"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App) {
	app.Post("/signup", validators.UserCreateValidator, auth_controllers.SignUp)
	app.Post("/login", auth_controllers.LogIn)
	app.Post("/refresh", auth_controllers.Refresh)

	// app.Post("/early_access", auth_controllers.GetEarlyAccessToken)

	app.Post("/recovery", auth_controllers.SendResetURL)
	app.Post("/recovery/verify", auth_controllers.ResetPassword)

	userRoutes := app.Group("/users", middlewares.Protect)
	userRoutes.Get("/me", user_controllers.GetMe)
	userRoutes.Get("/me/likes", user_controllers.GetMyLikes)
	userRoutes.Get("/me/organization/memberships", user_controllers.GetMyOrgMemberships)
	userRoutes.Get("/views", user_controllers.GetViews)

	userRoutes.Patch("/update_password", user_controllers.UpdatePassword)
	userRoutes.Patch("/update_email", user_controllers.UpdateEmail)
	userRoutes.Patch("/update_phone_number", user_controllers.UpdatePhoneNo)
	userRoutes.Patch("/update_resume", user_controllers.UpdateResume)

	userRoutes.Get("/deactivate", user_controllers.SendDeactivateVerificationCode)
	userRoutes.Post("/deactivate", user_controllers.Deactivate)

	userRoutes.Patch("/me", user_controllers.UpdateMe)
	userRoutes.Patch("/me/profile", user_controllers.EditProfile)
	userRoutes.Patch("/me/achievements", user_controllers.AddAchievement)
	userRoutes.Delete("/me/achievements/:achievementID", user_controllers.DeleteAchievement)
	// userRoutes.Delete("/me", controllers.DeactivateMe)

	userRoutes.Post("/report", controllers.AddReport)
	userRoutes.Post("/feedback", controllers.AddFeedback)
}
