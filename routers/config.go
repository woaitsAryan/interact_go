package routers

import "github.com/gofiber/fiber/v2"

func Config(app *fiber.App) {
	ConnectionRouter(app)
	PostRouter(app)
	ProjectRouter(app)
	UserRouter(app)
	FeedRouter(app)
	ApplicationRouter(app)
	BookmarkRouter(app)
	CommentRouter(app)
	ExploreRouter(app)
	InvitationRouter(app)
	MessagingRouter(app)
	NotificationRouter(app)
	OpeningRouter(app)
	WorkspaceRouter(app)
}
