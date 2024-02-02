package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/project_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ApplicationRouter(app *fiber.App) {
	applicationRoutes := app.Group("/applications", middlewares.Protect)

	applicationRoutes.Get("/:applicationID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.GetApplication("project"))

	applicationRoutes.Get("/accept/:applicationID", middlewares.ProjectRoleAuthorization(models.ProjectManager), project_controllers.AcceptApplication)
	applicationRoutes.Get("/reject/:applicationID", middlewares.ProjectRoleAuthorization(models.ProjectManager), project_controllers.RejectApplication)
	applicationRoutes.Get("/review/:applicationID", middlewares.ProjectRoleAuthorization(models.ProjectManager), project_controllers.SetApplicationReviewStatus)

	applicationRoutes.Post("/:openingID", controllers.AddApplication("project"))
	applicationRoutes.Delete("/:applicationID", controllers.DeleteApplication)
}
