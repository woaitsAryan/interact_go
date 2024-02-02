package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func OrgApplicationRouter(app *fiber.App) {
	applicationRoutes := app.Group("/org/:orgID/applications", middlewares.Protect)

	applicationRoutes.Get("/:applicationID", middlewares.OrgRoleAuthorization(models.Manager), controllers.GetApplication("org"))

	applicationRoutes.Get("/accept/:applicationID", middlewares.OrgRoleAuthorization(models.Manager), organization_controllers.AcceptApplication)
	applicationRoutes.Get("/reject/:applicationID", middlewares.OrgRoleAuthorization(models.Manager), organization_controllers.RejectApplication)
	applicationRoutes.Get("/review/:applicationID", middlewares.OrgRoleAuthorization(models.Manager), organization_controllers.SetApplicationReviewStatus)

	applicationRoutes.Post("/:openingID", controllers.AddApplication("org"))
}
