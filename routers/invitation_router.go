package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func InvitationRouter(app *fiber.App) {
	invitationRoutes := app.Group("/invitations", middlewares.Protect)
	invitationRoutes.Get("/me", controllers.GetInvitations)

	invitationRoutes.Get("/accept/chat/:invitationID", controllers.AcceptChatInvitation)
	invitationRoutes.Get("/accept/project/:invitationID", controllers.AcceptProjectInvitation)

	invitationRoutes.Get("/reject/chat/:invitationID", controllers.RejectChatInvitation)
	invitationRoutes.Get("/reject/project/:invitationID", controllers.RejectProjectInvitation)

	invitationRoutes.Get("/unread", controllers.GetUnreadInvitations)
	invitationRoutes.Post("/unread", controllers.MarkReadInvitations)

	invitationRoutes.Delete("/withdraw/chat/:invitationID", controllers.WithdrawChatInvitation)
	invitationRoutes.Delete("/withdraw/project/:invitationID", controllers.WithdrawProjectInvitation)
}
