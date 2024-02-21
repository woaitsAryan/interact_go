package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func MembershipRouter(app *fiber.App) {
	app.Get("/org/:orgID/membership/delete", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.SendLeaveOrgVerificationCode)
	app.Delete("/org/:orgID/membership", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.LeaveOrganization)

	app.Get("/org/:orgID/is_member", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
		})
	})

	app.Get("/org/:orgID/membership", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetMemberships)
	app.Get("/org/:orgID/explore_memberships", middlewares.PartialProtect, organization_controllers.GetExploreMemberships)

	membershipRoutes := app.Group("/org/:orgID/membership", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Manager))
	membershipRoutes.Get("/non_members", organization_controllers.GetNonMembers)
	membershipRoutes.Post("/", organization_controllers.AddMember)
	membershipRoutes.Patch("/:membershipID", organization_controllers.ChangeMemberRole)
	membershipRoutes.Delete("/:membershipID", organization_controllers.RemoveMember)
	membershipRoutes.Delete("/withdraw/:invitationID", controllers.WithdrawInvitation)
}
