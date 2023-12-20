package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/auth_controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func MiscRouter(app *fiber.App) {
	miscRouter := app.Group("/org/:orgID", middlewares.Protect)

	miscRouter.Get("/", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrganization)
	miscRouter.Patch("/", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.UpdateOrg)
	miscRouter.Patch("/profile", middlewares.OrgRoleAuthorization(models.Senior), controllers.EditProfile)
	miscRouter.Get("/history", middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrganizationHistory)

	miscRouter.Get("/delete", auth_controllers.SendDeleteVerificationCode)
	miscRouter.Delete("/delete", organization_controllers.DeleteOrganization)
}
