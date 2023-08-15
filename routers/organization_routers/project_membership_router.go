package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ProjectMembershipRouter(app *fiber.App) {
	membershipRoutes := app.Group("/org/project/membership", middlewares.Protect, middlewares.RoleAuthorization("Manager"))
	membershipRoutes.Post("/:projectID", controllers.AddMember)
	membershipRoutes.Patch("/:membershipID", controllers.ChangeMemberRole)
	membershipRoutes.Delete("/:projectID", controllers.LeaveProject)
	membershipRoutes.Delete("/:membershipID", controllers.RemoveMember)
}
