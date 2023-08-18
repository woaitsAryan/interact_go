package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func CommentRouter(app *fiber.App) {
	CommentRoutes := app.Group("/comments", middlewares.Protect)

	CommentRoutes.Get("/me/likes", controllers.GetMyLikedComments)

	CommentRoutes.Get("/post/:postID", controllers.GetPostComments)
	CommentRoutes.Get("/project/:projectID", controllers.GetProjectComments)

	CommentRoutes.Post("/", controllers.AddComment)

	CommentRoutes.Patch("/:commentID", controllers.UpdateComment)

	CommentRoutes.Delete("/:commentID", controllers.DeleteComment)

	CommentRoutes.Get("/like/:commentID", controllers.LikeComment)
}
