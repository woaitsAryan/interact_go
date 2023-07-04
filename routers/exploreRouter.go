package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ExploreRouter(app *fiber.App) {
	exploreRoutes := app.Group("/explore", middlewares.PartialProtect)
	exploreRoutes.Get("/posts", controllers.GetTrendingPosts)
	exploreRoutes.Get("/openings", controllers.GetTrendingOpenings)
	exploreRoutes.Get("/projects", controllers.GetTrendingProjects)
	exploreRoutes.Get("/users", controllers.GetTrendingUsers)
	exploreRoutes.Get("/users/:userID", controllers.GetUser)
	exploreRoutes.Get("/projects/:projectID", controllers.GetProject)
}
