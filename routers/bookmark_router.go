package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func BookmarkRouter(app *fiber.App) {
	bookmarkRoutes := app.Group("/bookmarks", middlewares.Protect)
	bookmarkRoutes.Get("/", controllers.GetBookMarks)

	bookmarkRoutes.Get("/post", controllers.GetPopulatedBookMarks("post"))
	bookmarkRoutes.Get("/project", controllers.GetPopulatedBookMarks("project"))
	bookmarkRoutes.Get("/opening", controllers.GetPopulatedBookMarks("opening"))

	bookmarkRoutes.Post("/post", controllers.AddBookMark("post"))
	bookmarkRoutes.Post("/project", controllers.AddBookMark("project"))
	bookmarkRoutes.Post("/opening", controllers.AddBookMark("opening"))

	bookmarkRoutes.Delete("/post/:bookmarkID", controllers.DeleteBookMark("post"))
	bookmarkRoutes.Delete("/project/:bookmarkID", controllers.DeleteBookMark("project"))
	bookmarkRoutes.Delete("/opening/:bookmarkID", controllers.DeleteBookMark("opening"))

	bookmarkRoutes.Post("/post/item/:bookmarkID", controllers.AddBookMarkItem("post"))
	bookmarkRoutes.Post("/project/item/:bookmarkID", controllers.AddBookMarkItem("project"))
	bookmarkRoutes.Post("/opening/item/:bookmarkID", controllers.AddBookMarkItem("opening"))

	bookmarkRoutes.Delete("/post/item/:bookmarkItemID", controllers.DeleteBookMarkItem("post"))
	bookmarkRoutes.Delete("/project/item/:bookmarkItemID", controllers.DeleteBookMarkItem("project"))
	bookmarkRoutes.Delete("/opening/item/:bookmarkItemID", controllers.DeleteBookMarkItem("opening"))
}
