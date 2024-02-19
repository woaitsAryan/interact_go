package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/controllers/user_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func MiscRouter(app *fiber.App) {
	miscRouter := app.Group("/org/:orgID", middlewares.Protect)

	miscRouter.Get("/", middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrganization)
	miscRouter.Patch("/", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.UpdateOrg)
	miscRouter.Patch("/me", middlewares.OrgRoleAuthorization(models.Senior), user_controllers.UpdateMe)
	miscRouter.Get("/history", middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrganizationHistory)

	miscRouter.Get("/newsfeed", organization_controllers.GetOrgNewsFeed)

	miscRouter.Get("/delete", user_controllers.SendDeactivateVerificationCode)
	miscRouter.Post("/delete", organization_controllers.DeleteOrganization)
}
