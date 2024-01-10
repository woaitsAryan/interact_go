package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ReviewRouter(app *fiber.App) {
	app.Get("/org/:orgID/review", organization_controllers.FetchOrgReviews)

	//TODO change all like/dislike routes to POST
	app.Get("/org/:orgID/review/like/:reviewID", middlewares.Protect, controllers.LikeItem("review"))
	app.Get("/org/:orgID/review/dislike/:reviewID", middlewares.Protect, controllers.DislikeItem("review"))

	reviewRoutes := app.Group("/org/:orgID/review", middlewares.OrgRoleAuthorization(models.Member), middlewares.Protect)
	reviewRoutes.Post("/", organization_controllers.AddReview)
	reviewRoutes.Delete("/", organization_controllers.DeleteReview)
}
