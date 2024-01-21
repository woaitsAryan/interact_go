package organization_routers

import (
	"github.com/gofiber/fiber/v2"
)

func Config(app *fiber.App) {
	AnnouncementRouter(app)
	AuthRouter(app)
	ChatRouter(app)
	EventRouter(app)
	MembershipRouter(app)
	MiscRouter(app)
	PollRouter(app)
	PostRouter(app)
	ProjectApplicationRouter(app)
	ProjectMembershipRouter(app)
	ProjectOpeningRouter(app)
	ProjectRouter(app)
	ResourceRouter(app)
	ReviewRouter(app)
	TaskRouter(app)
	OrgOpeningRouter(app)
	OrgApplicationRouter(app)
}
