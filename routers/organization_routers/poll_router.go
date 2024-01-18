package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

/* Router for adding polls */
func PollRouter(app *fiber.App) {
	pollRouter := app.Group("/org/:orgID/polls", middlewares.Protect)
	pollRouter.Get("/", organization_controllers.FetchPolls)

	pollRouter.Post("/", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.CreatePoll)
	pollRouter.Patch("/:pollID", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.EditPoll)
	pollRouter.Delete("/:pollID", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.DeletePoll)

	pollRouter.Patch("/vote/:pollID/:OptionID", middlewares.OrgPollAuthorization(), organization_controllers.VotePoll)
	pollRouter.Patch("/unvote/:pollID/:OptionID", middlewares.OrgPollAuthorization(), organization_controllers.UnvotePoll)
}
