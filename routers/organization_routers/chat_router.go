package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/messaging_controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ChatRouter(app *fiber.App) {

	app.Get("/org/:orgID/chats", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrganizationChats)

	chatRoutes := app.Group("/org/:orgID/chats", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	//TODO36 add chat and chat membership edit routes for org acc and org managers
	chatRoutes.Post("/", messaging_controllers.AddGroupChat("Organization"))
}
