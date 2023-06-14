package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ConnectionRouter(app *fiber.App) {
	app.Get("/follow/:userID", middlewares.Protect, controllers.FollowUser)
	app.Get("/unfollow/:userID", middlewares.Protect, controllers.UnfollowUser)
	app.Get("/remove_follow/:userID", middlewares.Protect, controllers.RemoveFollow)
	app.Get("/followers/:userID", middlewares.Protect, controllers.GetFollowers)
	app.Get("/following/:userID", middlewares.Protect, controllers.GetFollowing)
}
