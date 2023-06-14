package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func PostRouter(app *fiber.App) {
	postRoutes := app.Group("/posts", middlewares.Protect)
	postRoutes.Post("/", controllers.AddPost)
	postRoutes.Get("/me", controllers.GetMyPosts)
	postRoutes.Get("/:postID", controllers.GetPost)
	postRoutes.Patch("/:postID", middlewares.PostUserProtect, controllers.UpdatePost)
	postRoutes.Delete("/:postID", middlewares.PostUserProtect, controllers.DeletePost)
}
