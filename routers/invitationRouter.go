package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func InvitationRouter(app *fiber.App) {
	exploreRoutes := app.Group("/invitations", middlewares.Protect)
	exploreRoutes.Get("/me", controllers.GetInvitations)

	exploreRoutes.Get("/accept/chat/:invitationID", controllers.AcceptChatInvitation)
	exploreRoutes.Get("/accept/project/:invitationID", controllers.AcceptProjectInvitation)

	exploreRoutes.Get("/reject/chat/:invitationID", controllers.RejectChatInvitation)
	exploreRoutes.Get("/reject/project/:invitationID", controllers.RejectProjectInvitation)

	exploreRoutes.Delete("/withdraw/chat/:invitationID", controllers.WithdrawChatInvitation)
	exploreRoutes.Delete("/withdraw/project/:invitationID", controllers.WithdrawProjectInvitation)
}
