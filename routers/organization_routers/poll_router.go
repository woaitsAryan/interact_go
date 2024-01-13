package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

/* Router for adding polls */
func PollRouter(app *fiber.App) {
	pollRouter := app.Group("/org/:orgID/poll", middlewares.Protect)
	pollRouter.Post("/", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.CreatePoll)
	pollRouter.Get("/", middlewares.OrgRoleAuthorization(models.Member), organization_controllers.FetchPolls)
	pollRouter.Patch("/vote/:pollID/:OptionID", middlewares.OrgRoleAuthorization(models.Member),  organization_controllers.VotePoll)
	pollRouter.Patch("/unvote/:OptionID", middlewares.OrgRoleAuthorization(models.Member), organization_controllers.UnvotePoll)
	pollRouter.Delete("/:pollID", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.DeletePoll)
	pollRouter.Patch("/:pollID", middlewares.OrgRoleAuthorization(models.Senior), organization_controllers.EditPoll)
}
