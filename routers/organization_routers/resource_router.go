package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ResourceRouter(app *fiber.App) {
	resourceRoutes := app.Group("/org/:orgID/resource", middlewares.Protect)

	resourceRoutes.Get("/", middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrgResourceBuckets)
	resourceRoutes.Get("/:resourceBucketID", middlewares.OrgBucketAuthorization("view"), organization_controllers.GetResourceBucketFiles)
	resourceRoutes.Post("/", middlewares.OrgRoleAuthorization(models.Manager), organization_controllers.AddResourceBucket)
	resourceRoutes.Patch("/:resourceBucketID", middlewares.OrgBucketAuthorization("edit"), organization_controllers.EditResourceBucket)
	resourceRoutes.Delete("/:resourceBucketID", middlewares.OrgRoleAuthorization(models.Manager), organization_controllers.DeleteResourceBucket)

	resourceFileRoutes := resourceRoutes.Group("/:resourceBucketID/file")

	resourceFileRoutes.Post("/", middlewares.OrgBucketAuthorization("edit"), organization_controllers.AddResourceBucket)
	resourceFileRoutes.Patch("/:resourceFileID", middlewares.OrgBucketAuthorization("edit"), organization_controllers.EditResourceBucket)
	resourceFileRoutes.Delete("/:resourceFileID", middlewares.OrgBucketAuthorization("edit"), organization_controllers.DeleteResourceBucket)
}
