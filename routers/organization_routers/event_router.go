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

	eventRoutesOrg := app.Group("/org/:orgID/events", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	eventRoutesOrg.Post("/",organization_controllers.AddEvent)
	eventRoutesOrg.Delete("/:eventID", organization_controllers.DeleteEvent)
	eventRoutesOrg.Post("/cohost", organization_controllers.AddOtherOrg)
	eventRoutesOrg.Delete("/cohost", organization_controllers.RemoveOtherOrg)

	eventRoutesOrgCoOwn := app.Group("/org/:orgID/events", middlewares.Protect, middlewares.OrgEventRoleAuthorization(models.Senior))
	eventRoutesOrgCoOwn.Post("/coordinators/:eventID", organization_controllers.AddEventCoordinators)
	eventRoutesOrgCoOwn.Delete("/coordinators/:eventID", organization_controllers.RemoveEventCoordinators)
	eventRoutesOrgCoOwn.Patch("/:eventID",organization_controllers.UpdateEvent)
	
}
