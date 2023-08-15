package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MembershipRouter(app *fiber.App) {
	membershipRoutes := app.Group("/org/membership", middlewares.Protect, middlewares.RoleAuthorization("Owner"))
	membershipRoutes.Post("/project/:projectID", controllers.AddMember)
	membershipRoutes.Patch("/:membershipID", controllers.ChangeMemberRole)
	membershipRoutes.Delete("project/:projectID", controllers.LeaveProject) // different middlewareAuthorization
	membershipRoutes.Delete("/:membershipID", controllers.RemoveMember)
}
