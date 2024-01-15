package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func AnnouncementRouter(app *fiber.App) {
	app.Get("/org/:orgID/announcements", middlewares.Protect, organization_controllers.GetOrgAnnouncements)

	announcementRoutes := app.Group("/org/:orgID/announcements", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	announcementRoutes.Post("/", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.AddResourceBucket)
	announcementRoutes.Patch("/:announcementID", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.EditResourceBucket)
	announcementRoutes.Delete("/:announcementID", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.DeleteResourceBucket)
}
