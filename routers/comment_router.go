package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func CommentRouter(app *fiber.App) {
	commentRoutes := app.Group("/comments", middlewares.Protect)

	commentRoutes.Get("/post/:postID", controllers.GetPostComments)
	commentRoutes.Get("/project/:projectID", controllers.GetProjectComments)
	commentRoutes.Get("/event/:eventID", controllers.GetEventComments)

	commentRoutes.Post("/", controllers.AddComment)

	commentRoutes.Patch("/:commentID", controllers.UpdateComment)

	commentRoutes.Delete("/:commentID", controllers.DeleteComment)

	commentRoutes.Get("/like/:commentID", controllers.LikeItem("comment"))
	commentRoutes.Get("/dislike/:commentID", controllers.DislikeItem("comment"))
}
