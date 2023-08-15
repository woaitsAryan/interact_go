package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func PostRouter(app *fiber.App) {
	postRoutes := app.Group("/org/posts", middlewares.Protect, middlewares.RoleAuthorization(models.Manager))
	postRoutes.Post("/", controllers.AddPost)
	postRoutes.Patch("/:postID", controllers.UpdatePost)
	postRoutes.Delete("/:postID", controllers.DeletePost)
}
