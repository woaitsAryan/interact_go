package routers

import (
	"github.com/Pratham-Mishra04/interact/routers/organization_routers"
	"github.com/gofiber/fiber/v2"
)

func Config(app *fiber.App) {
	UserRouter(app)
	OauthRouter(app)

	ConnectionRouter(app)
	PostRouter(app)
	ProjectRouter(app)

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
	MembershipRouter(app)
	ShareRouter(app)
	TaskRouter(app)

	VerificationRouter(app)

	organization_routers.Config(app)
}
