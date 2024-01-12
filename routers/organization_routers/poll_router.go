package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

/* Router for adding polls */
func PollRouter(app *fiber.App) {
	pollRouter := app.Group("/org/:orgID/poll", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member))
	pollRouter.Post("/", organization_controllers.CreatePoll)
	pollRouter.Get("/", organization_controllers.FetchPolls)
	pollRouter.Patch("/vote/:pollID/:OptionID", organization_controllers.VotePoll)
	pollRouter.Patch("/unvote/:OptionID", organization_controllers.UnvotePoll)
	pollRouter.Delete("/:pollID", organization_controllers.DeletePoll)
}
