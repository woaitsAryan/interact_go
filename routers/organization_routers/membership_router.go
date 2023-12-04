package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func MembershipRouter(app *fiber.App) {

	app.Delete("/org/:orgID/membership", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.LeaveOrganization)

	app.Get("/org/:orgID/membership", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetMemberships)

	membershipRoutes := app.Group("/org/:orgID/membership", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Manager))
	membershipRoutes.Get("/non_members", organization_controllers.GetNonMembers)
	membershipRoutes.Post("/", organization_controllers.AddMember)
	membershipRoutes.Patch("/:membershipID", organization_controllers.ChangeMemberRole)
	membershipRoutes.Delete("/:membershipID", organization_controllers.RemoveMember)
}
