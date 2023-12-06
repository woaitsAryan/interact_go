package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func Config(app *fiber.App) {
	app.Get("/org/:orgID", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrganization)

	AuthRouter(app)
	ChatRouter(app)
	PostRouter(app)
	ProjectRouter(app)
	ProjectApplicationRouter(app)
	ProjectMembershipRouter(app)
	ProjectOpeningRouter(app)
	MembershipRouter(app)
	TaskRouter(app)
}
