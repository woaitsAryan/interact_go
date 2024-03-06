package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

/*
	Router for reviewing organizations.

Add and delete a review.
Like, dislike, removing that like and dislike from that review.
Fetch all reviews.
*/
func ReviewRouter(app *fiber.App) {
	app.Get("/org/:orgID/reviews", organization_controllers.FetchOrgReviews)
	app.Get("/org/:orgID/reviews/data", organization_controllers.GetOrgReviewData)

	//TODO38 change all like/dislike routes to POST
	app.Get("/org/:orgID/reviews/like/:reviewID", middlewares.Protect, controllers.LikeItem("review"))
	app.Get("/org/:orgID/reviews/dislike/:reviewID", middlewares.Protect, controllers.DislikeItem("review"))

	reviewRoutes := app.Group("/org/:orgID/reviews", middlewares.OrgRoleAuthorization(models.Member), middlewares.Protect)
	reviewRoutes.Post("/", organization_controllers.AddReview)
	reviewRoutes.Delete("/:reviewID", organization_controllers.DeleteReview)
}
