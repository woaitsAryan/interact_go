package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func EventRouter(app *fiber.App) {
	app.Get("/events/like/:eventID", middlewares.Protect, controllers.LikeItem("event"))
	app.Get("/events/dislike/:eventID", middlewares.Protect, controllers.DislikeItem("event"))

	eventRoutesOrg := app.Group("/org/:orgID/events", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	eventRoutesOrg.Get("/invitations", controllers.GetInvitations)
	eventRoutesOrg.Post("/", organization_controllers.AddEvent)
	eventRoutesOrg.Delete("/:eventID", organization_controllers.DeleteEvent)

	eventRoutesOrg.Post("/:eventID/cohost", organization_controllers.AddCoHostOrgs)
	eventRoutesOrg.Delete("/:eventID/cohost", organization_controllers.RemoveCoHostOrg)

	eventRoutesOrg.Patch("/:eventID", middlewares.OrgEventCoHostAuthorization, organization_controllers.UpdateEvent)
	eventRoutesOrg.Patch("/:eventID/cohost", middlewares.OrgEventCoHostAuthorization, organization_controllers.LeaveCoHostOrg)

	eventRoutesOrg.Post("/coordinators/:eventID", middlewares.OrgEventCoHostAuthorization, organization_controllers.AddEventCoordinators)
	eventRoutesOrg.Delete("/coordinators/:eventID", middlewares.OrgEventCoHostAuthorization, organization_controllers.RemoveEventCoordinators)
}
