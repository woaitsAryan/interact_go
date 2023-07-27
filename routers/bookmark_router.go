package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func BookmarkRouter(app *fiber.App) {
	bookmarkRoutes := app.Group("/bookmarks", middlewares.Protect)
	bookmarkRoutes.Get("/", controllers.GetBookMarks)

	bookmarkRoutes.Get("/post", controllers.GetPopulatedPostBookMarks)
	bookmarkRoutes.Get("/project", controllers.GetPopulatedProjectBookMarks)

	bookmarkRoutes.Post("/post", controllers.AddPostBookMark)
	bookmarkRoutes.Post("/project", controllers.AddProjectBookMark)

	bookmarkRoutes.Delete("/post/:bookmarkID", controllers.DeletePostBookMark)
	bookmarkRoutes.Delete("/project/:bookmarkID", controllers.DeleteProjectBookMark)

	bookmarkRoutes.Post("/post/item/:bookmarkID", controllers.AddPostBookMarkItem)
	bookmarkRoutes.Post("/project/item/:bookmarkID", controllers.AddProjectBookMarkItem)

	bookmarkRoutes.Delete("/post/item/:bookmarkItemID", controllers.DeletePostBookMarkItem)
	bookmarkRoutes.Delete("/project/item/:bookmarkItemID", controllers.DeleteProjectBookMarkItem)
}
