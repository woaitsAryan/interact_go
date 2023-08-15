package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func MembershipRouter(app *fiber.App) {

	app.Delete("/org/membership/:organizationID", middlewares.Protect, middlewares.RoleAuthorization(models.Member), organization_controllers.LeaveOrganization)

	membershipRoutes := app.Group("/org/membership", middlewares.Protect, middlewares.RoleAuthorization(models.Owner))
	membershipRoutes.Post("/:organizationID", organization_controllers.AddMember)
	membershipRoutes.Patch("/:membershipID", organization_controllers.ChangeMemberRole)
	membershipRoutes.Delete("/:membershipID", organization_controllers.RemoveMember)
}
