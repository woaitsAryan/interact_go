package organization_routers

import "github.com/gofiber/fiber/v2"

func Config(app *fiber.App) {
	OauthRouter(app)
	PostRouter(app)
	ProjectRouter(app)
	ProjectApplicationRouter(app)
	CommentRouter(app)
	ProjectMembershipRouter(app)
	ProjectOpeningRouter(app)
}
