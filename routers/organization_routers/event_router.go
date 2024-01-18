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
	app.Get("/events/like/:eventID", middlewares.Protect, controllers.LikeItem("event"))
	app.Get("/events/dislike/:eventID", middlewares.Protect, controllers.DislikeItem("event"))

	eventRoutes := app.Group("/org/:orgID/events", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	eventRoutes.Post("/", organization_controllers.AddEvent)

	eventRoutes.Post("/coordinators/:eventID", organization_controllers.AddEventCoordinators)
	eventRoutes.Delete("/coordinators/:eventID", organization_controllers.RemoveEventCoordinators)

	eventRoutes.Patch("/:eventID", organization_controllers.UpdateEvent)
	eventRoutes.Delete("/:eventID", organization_controllers.DeleteEvent)
}
