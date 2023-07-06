package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func CommentRouter(app *fiber.App) {
	CommentRoutes := app.Group("/comments", middlewares.Protect)

	CommentRoutes.Get("/post/me/likes", controllers.GetMyLikedPostsComments)
	CommentRoutes.Get("/project/me/likes", controllers.GetMyLikedProjectsComments)

	CommentRoutes.Get("/post/:postID", controllers.GetPostComments)
	CommentRoutes.Get("/project/:projectID", controllers.GetProjectComments)

	CommentRoutes.Post("/post", controllers.AddPostComment)
	CommentRoutes.Post("/project", controllers.AddProjectComment)

	CommentRoutes.Patch("/:commentID", controllers.UpdatePostComment)

	CommentRoutes.Delete("post/:commentID", controllers.DeletePostComment)
	CommentRoutes.Delete("project/:commentID", controllers.DeleteProjectComment)

	CommentRoutes.Get("/post/like/:commentID", controllers.LikePostComment)
	CommentRoutes.Get("/project/like/:commentID", controllers.LikeProjectComment)
}
