package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func BookmarkRouter(app *fiber.App) {
	bookmarkRoutes := app.Group("/bookmarks", middlewares.Protect)
	bookmarkRoutes.Get("/", controllers.GetBookMarks)

	bookmarkRoutes.Post("/post", controllers.AddPostBookMark)
	bookmarkRoutes.Post("/project", controllers.AddProjectBookMark)

	bookmarkRoutes.Delete("/post/:bookmarkID", controllers.DeletePostBookMark)
	bookmarkRoutes.Delete("/project/:bookmarkID", controllers.DeleteProjectBookMark)

	bookmarkRoutes.Post("/post/item", controllers.AddPostBookMarkItem)
	bookmarkRoutes.Post("/project/item", controllers.AddProjectBookMarkItem)

	bookmarkRoutes.Delete("/post/item/:bookmarkID", controllers.DeletePostBookMarkItem)
	bookmarkRoutes.Delete("/project/item/:bookmarkID", controllers.DeleteProjectBookMarkItem)
}
