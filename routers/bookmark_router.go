package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func setupBookmarkRoutes(router fiber.Router, bookmarkType string) {
	router.Get("/"+bookmarkType, controllers.GetPopulatedBookMarks(bookmarkType))
	router.Post("/"+bookmarkType, controllers.AddBookMark(bookmarkType))
	router.Patch("/"+bookmarkType+"/:bookmarkID", controllers.UpdateBookMark(bookmarkType))
	router.Delete("/"+bookmarkType+"/:bookmarkID", controllers.DeleteBookMark(bookmarkType))

	router.Post("/"+bookmarkType+"/item/:bookmarkID", controllers.AddBookMarkItem(bookmarkType))
	router.Delete("/"+bookmarkType+"/item/:bookmarkItemID", controllers.DeleteBookMarkItem(bookmarkType))
}

func BookmarkRouter(app *fiber.App) {
	bookmarkRouter := app.Group("/bookmarks", middlewares.Protect)

	setupBookmarkRoutes(bookmarkRouter, "post")
	setupBookmarkRoutes(bookmarkRouter, "project")
	setupBookmarkRoutes(bookmarkRouter, "opening")
	setupBookmarkRoutes(bookmarkRouter, "event")
}
