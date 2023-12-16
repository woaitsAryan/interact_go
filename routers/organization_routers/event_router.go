package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func EventRouter(app *fiber.App) {
	app.Get("/org/:orgID/events", middlewares.Protect, organization_controllers.GetOrgEvents)
	app.Patch("events/like/:eventID", middlewares.Protect, controllers.LikeEvent)

	eventRoutes := app.Group("/org/:orgID/events", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	eventRoutes.Post("/", organization_controllers.AddEvent)
	eventRoutes.Patch("/:eventID", organization_controllers.UpdateEvent)
	eventRoutes.Delete("/:eventID", organization_controllers.DeleteEvent)
}
