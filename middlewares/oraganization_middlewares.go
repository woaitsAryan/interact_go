package middlewares

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func checkAccess(UserRole models.OrganizationRole, AuthorizedRole models.OrganizationRole) bool {
	if UserRole == models.Owner {
		return true
	} else if UserRole == models.Manager {
		return AuthorizedRole != models.Owner
	} else if UserRole == models.Member {
		return AuthorizedRole == models.Member
	}

	return false
}

func RoleAuthorization(Role models.OrganizationRole) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		orgID := c.Params("orgID")

		var orgMembership models.OrganizationMembership
		if err := initializers.DB.Preload("Organization").First(orgMembership, "organization_id = ? AND user_id=?", orgID, loggedInUserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				var org models.Organization
				if err := initializers.DB.First(org, "user_id=?", loggedInUserID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return &fiber.Error{Code: 403, Message: "Cannot access this organization"}
					}
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
				}
				return c.Next()
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}

		if !checkAccess(orgMembership.Role, Role) {
			return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
		}

		c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
		c.Set("loggedInUserID", orgMembership.Organization.UserID.String())

		return c.Next()
	}
}
