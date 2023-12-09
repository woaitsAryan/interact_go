package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func MiscRouter(app *fiber.App) {
	miscRouter := app.Group("/org/:orgID", middlewares.Protect)

	miscRouter.Get("/history", middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrganizationHistory)
}
