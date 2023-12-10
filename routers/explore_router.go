package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/explore_controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/controllers/project_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ExploreRouter(app *fiber.App) {
	exploreRoutes := app.Group("/explore", middlewares.PartialProtect)

	exploreRoutes.Get("/trending_searches", explore_controllers.GetTrendingSearches)
	exploreRoutes.Post("/search", explore_controllers.AddSearchQuery)

	exploreRoutes.Get("/posts/trending", explore_controllers.GetTrendingPosts)
	exploreRoutes.Get("/posts/latest", explore_controllers.GetLatestPosts)
	exploreRoutes.Get("/posts/recommended", explore_controllers.GetRecommendedPosts)

	exploreRoutes.Get("/openings/recommended", explore_controllers.GetRecommendedOpenings)
	exploreRoutes.Get("/openings/trending", explore_controllers.GetTrendingOpenings)
	exploreRoutes.Get("/openings/:slug", explore_controllers.GetProjectOpenings)

	exploreRoutes.Get("/projects/trending", explore_controllers.GetTrendingProjects)
	exploreRoutes.Get("/projects/recommended", explore_controllers.GetRecommendedProjects)

	exploreRoutes.Get("/events/trending", explore_controllers.GetTrendingEvents)
	exploreRoutes.Get("/events/recommended", explore_controllers.GetRecommendedEvents)
	exploreRoutes.Get("/events/org/:orgID", explore_controllers.GetOrgEvents)
	exploreRoutes.Get("/events/:eventID", organization_controllers.GetEvent)

	exploreRoutes.Get("/projects/most_liked", explore_controllers.GetMostLikedProjects)
	exploreRoutes.Get("/projects/recently_added", middlewares.Protect, explore_controllers.GetLatestProjects)
	exploreRoutes.Get("/projects/last_viewed", middlewares.Protect, explore_controllers.GetLastViewedProjects)

	exploreRoutes.Get("/users/trending", explore_controllers.GetTrendingUsers)
	exploreRoutes.Get("/users/recommended", explore_controllers.GetRecommendedUsers)

	exploreRoutes.Get("/orgs/trending", explore_controllers.GetTrendingOrganizationalUsers)
	exploreRoutes.Get("/orgs/recommended", explore_controllers.GetRecommendedOrganizationalUsers)
	exploreRoutes.Get("/orgs/:username", explore_controllers.GetOrganizationalUser)

	exploreRoutes.Get("/users/similar/:username", explore_controllers.GetSimilarUsers)
	exploreRoutes.Get("/projects/similar/:slug", explore_controllers.GetSimilarProjects)
	exploreRoutes.Get("/events/similar/:eventID", explore_controllers.GetSimilarEvents)

	exploreRoutes.Get("/users/posts/:userID", controllers.GetUserPosts)
	exploreRoutes.Get("/users/projects/:userID", project_controllers.GetUserProjects)
	exploreRoutes.Get("/users/projects/contributing/:userID", project_controllers.GetUserContributingProjects)

	exploreRoutes.Get("/users/:username", controllers.GetUser)
	exploreRoutes.Get("/projects/:slug", project_controllers.GetProject)
}
