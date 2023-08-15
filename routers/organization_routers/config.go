package organization_routers

import "github.com/gofiber/fiber/v2"

func Config(app *fiber.App) {
	AuthRouter(app)
	PostRouter(app)
	ProjectRouter(app)
	ProjectApplicationRouter(app)
	ProjectMembershipRouter(app)
	ProjectOpeningRouter(app)
}
