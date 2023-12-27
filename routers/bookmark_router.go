package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func setupBookmarkRoutes(app *fiber.App, bookmarkType string) {
	bookmarkRoutes := app.Group("/bookmarks", middlewares.Protect)

	bookmarkRoutes.Get("/"+bookmarkType, controllers.GetPopulatedBookMarks(bookmarkType))
	bookmarkRoutes.Post("/"+bookmarkType, controllers.AddBookMark(bookmarkType))
	bookmarkRoutes.Patch("/"+bookmarkType+"/:bookmarkID", controllers.UpdateBookMark(bookmarkType))
	bookmarkRoutes.Delete("/"+bookmarkType+"/:bookmarkID", controllers.DeleteBookMark(bookmarkType))
	bookmarkRoutes.Post("/"+bookmarkType+"/item/:bookmarkID", controllers.AddBookMarkItem(bookmarkType))
	bookmarkRoutes.Delete("/"+bookmarkType+"/item/:bookmarkItemID", controllers.DeleteBookMarkItem(bookmarkType))
}

func BookmarkRouter(app *fiber.App) {
	setupBookmarkRoutes(app, "post")
	setupBookmarkRoutes(app, "project")
	setupBookmarkRoutes(app, "opening")
	setupBookmarkRoutes(app, "event")
}
