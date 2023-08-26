package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ConnectionRouter(app *fiber.App) {
	connectionRoutes := app.Group("/connection", middlewares.Protect)

	connectionRoutes.Get("/follow/:userID", controllers.FollowUser)
	connectionRoutes.Get("/unfollow/:userID", controllers.UnfollowUser)
	connectionRoutes.Get("/remove_follow/:userID", controllers.RemoveFollow)
	connectionRoutes.Get("/followers/:userID", controllers.GetFollowers)
	connectionRoutes.Get("/following/:userID", controllers.GetFollowing)
	connectionRoutes.Get("/mutuals/:userID", controllers.GetMutuals)
}
