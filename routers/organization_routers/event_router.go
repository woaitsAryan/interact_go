package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers/event_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func EventRouter(app *fiber.App) {
	app.Get("/events/like/:eventID", middlewares.Protect, controllers.LikeItem("event"))
	app.Get("/events/dislike/:eventID", middlewares.Protect, controllers.DislikeItem("event"))

	eventRoutesOrg := app.Group("/org/:orgID/events", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	eventRoutesOrg.Get("/", event_controllers.GetPopulatedOrgEvents)
	eventRoutesOrg.Get("/invitations", controllers.GetInvitations)
	eventRoutesOrg.Get("/invitations/count", controllers.GetUnreadInvitationCount)
	eventRoutesOrg.Delete("/invitations/:invitationID", controllers.WithdrawInvitation)
	eventRoutesOrg.Post("/", event_controllers.AddEvent)
	eventRoutesOrg.Delete("/:eventID", event_controllers.DeleteEvent)

	eventRoutesOrg.Get("/:eventID/cohost", event_controllers.GetEventCoHosts)
	eventRoutesOrg.Post("/:eventID/cohost", event_controllers.AddCoHostOrgs)
	eventRoutesOrg.Delete("/:eventID/cohost", event_controllers.RemoveCoHostOrgs)

	eventRoutesOrg.Patch("/:eventID", middlewares.OrgEventCoHostAuthorization, event_controllers.UpdateEvent)
	eventRoutesOrg.Patch("/:eventID/cohost", middlewares.OrgEventCoHostAuthorization, event_controllers.LeaveCoHostOrg)

	eventRoutesOrg.Post("/coordinators/:eventID", middlewares.OrgEventCoHostAuthorization, event_controllers.AddEventCoordinators)
	eventRoutesOrg.Delete("/coordinators/:eventID", middlewares.OrgEventCoHostAuthorization, event_controllers.RemoveEventCoordinators)
}
