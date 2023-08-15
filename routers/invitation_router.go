package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func InvitationRouter(app *fiber.App) {
	invitationRoutes := app.Group("/invitations", middlewares.Protect)
	invitationRoutes.Get("/me", controllers.GetInvitations)

	invitationRoutes.Get("/accept/:invitationID", controllers.AcceptInvitation)

	invitationRoutes.Get("/reject/:invitationID", controllers.RejectInvitation)

	invitationRoutes.Get("/unread", controllers.GetUnreadInvitations)

	invitationRoutes.Post("/unread", controllers.MarkReadInvitations)

	invitationRoutes.Delete("/withdraw/:invitationID", controllers.WithdrawInvitation)
}
