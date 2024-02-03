package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func OrgOpeningRouter(app *fiber.App) {

	app.Get("/org_openings/:openingID", middlewares.PartialProtect, organization_controllers.GetOpening)

	orgOpeningRouter := app.Group("/org/:orgID/org_openings", middlewares.Protect)
	orgOpeningRouter.Get("/", organization_controllers.GetAllOpeningsOfOrganization)

	orgOpeningRouter.Get("/applications/:openingID", middlewares.OrgRoleAuthorization(models.Manager), controllers.GetAllApplicationsOfOpening)

	orgOpeningRouter.Post("/", middlewares.OrgRoleAuthorization(models.Manager), organization_controllers.AddOpening)
	orgOpeningRouter.Patch("/:openingID", middlewares.OrgRoleAuthorization(models.Manager), organization_controllers.EditOpening)
	orgOpeningRouter.Delete("/:openingID", middlewares.OrgRoleAuthorization(models.Manager), organization_controllers.DeleteOpening)
}
