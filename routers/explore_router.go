package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ExploreRouter(app *fiber.App) {
	exploreRoutes := app.Group("/explore")

	exploreRoutes.Get("/trending_searches", controllers.GetTrendingSearches)
	exploreRoutes.Post("/search", controllers.AddSearchQuery)

	exploreRoutes.Get("/posts", controllers.GetTrendingPosts)

	exploreRoutes.Get("/openings/recommended", controllers.GetRecommendedOpenings)
	exploreRoutes.Get("/openings/trending", controllers.GetTrendingOpenings)
	exploreRoutes.Get("/openings/:slug", controllers.GetProjectOpenings)

	exploreRoutes.Get("/projects/trending", controllers.GetTrendingProjects)
	exploreRoutes.Get("/projects/recommended", middlewares.PartialProtect, controllers.GetRecommendedProjects)
	exploreRoutes.Get("/projects/most_liked", controllers.GetMostLikedProjects)
	exploreRoutes.Get("/projects/recently_added", middlewares.Protect, controllers.GetRecentlyAddedProjects)
	exploreRoutes.Get("/projects/last_viewed", middlewares.Protect, controllers.GetLastViewedProjects)

	exploreRoutes.Get("/users/trending", controllers.GetTrendingUsers)
	exploreRoutes.Get("/users/recommended", middlewares.PartialProtect, controllers.GetRecommendedUsers)

	exploreRoutes.Get("/users/similar/:username", controllers.GetSimilarUsers)
	exploreRoutes.Get("/projects/similar/:slug", controllers.GetSimilarProjects)

	exploreRoutes.Get("/users/:username", controllers.GetUser)
	exploreRoutes.Get("/projects/:slug", controllers.GetProject)
}
