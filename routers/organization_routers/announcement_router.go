package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func AnnouncementRouter(app *fiber.App) {
	app.Get("/org/:orgID/announcements", middlewares.Protect, organization_controllers.GetOrgAnnouncements)
	app.Get("/org/:orgID/announcements/like/:announcementID", middlewares.Protect, controllers.LikeItem("announcement"))
	app.Get("/org/:orgID/announcements/dislike/:announcementID", middlewares.Protect, controllers.DislikeItem("announcement"))

	announcementRoutes := app.Group("/org/:orgID/announcements", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	announcementRoutes.Post("/", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.AddAnnouncement)
	announcementRoutes.Patch("/:announcementID", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.EditAnnouncement)
	announcementRoutes.Delete("/:announcementID", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.DeleteAnnouncement)
}
