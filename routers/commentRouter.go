package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func CommentRouter(app *fiber.App) {
	CommentRoutes := app.Group("/comments", middlewares.Protect)
	CommentRoutes.Get("/post/:postID", controllers.GetPostComments)
	// CommentRoutes.Get("/project/:projectID", controllers.GetProjectComments)
	CommentRoutes.Post("/post", controllers.AddPostComment)
	CommentRoutes.Patch("/:commentID", controllers.UpdatePostComment)
	CommentRoutes.Delete("post/:commentID", controllers.DeletePostComment)

	CommentRoutes.Get("/like/:commentID", controllers.LikePostComment)
}
