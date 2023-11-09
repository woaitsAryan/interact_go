package middlewares

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func checkOrgAccess(UserRole models.OrganizationRole, AuthorizedRole models.OrganizationRole) bool {
	if UserRole == models.Owner {
		return true
	} else if UserRole == models.Manager {
		return AuthorizedRole != models.Owner
	} else if UserRole == models.Member {
		return AuthorizedRole == models.Member
	}

	return false
}

func checkProjectAccess(UserRole models.ProjectRole, AuthorizedRole models.ProjectRole) bool {
	if UserRole == models.ProjectManager {
		return true
	} else if UserRole == models.ProjectEditor {
		return AuthorizedRole != models.ProjectManager
	} else if UserRole == models.ProjectMember {
		return AuthorizedRole == models.ProjectMember
	}

	return false
}

func OrgRoleAuthorization(Role models.OrganizationRole) func(*fiber.Ctx) error {
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

		if !checkOrgAccess(orgMembership.Role, Role) {
			return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
		}

		c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
		c.Set("loggedInUserID", orgMembership.Organization.UserID.String())

		return c.Next()
	}
}

func ProjectRoleAuthorization(Role models.ProjectRole) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		slug := c.Params("slug")
		projectID := c.Params("projectID")
		openingID := c.Params("openingID")
		applicationID := c.Params("applicationID")
		chatID := c.Params("chatID")
		membershipID := c.Params("membershipID")
		taskID := c.Params("taskID")

		var project models.Project
		if slug != "" {
			if err := initializers.DB.Preload("Memberships").First(&project, "slug = ?", slug).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid Project."}
			}
		} else if projectID != "" {
			if err := initializers.DB.Preload("Memberships").First(&project, "id = ?", projectID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid Project."}
			}
		} else if openingID != "" {
			var opening models.Opening
			if err := initializers.DB.Preload("Project").Preload("Project.Memberships").First(&opening, "id = ?", openingID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid Project."}
			}
			project = opening.Project
		} else if applicationID != "" {
			var application models.Application
			if err := initializers.DB.Preload("Project").Preload("Project.Memberships").First(&application, "id = ?", applicationID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid Project."}
			}
			project = application.Project
		} else if chatID != "" {
			var chat models.GroupChat
			if err := initializers.DB.Preload("Project").Preload("Project.Memberships").First(&chat, "id = ?", chatID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid Project."}
			}
			project = chat.Project
		} else if membershipID != "" {
			var membership models.Membership
			if err := initializers.DB.Preload("Project").Preload("Project.Memberships").First(&membership, "id = ?", membershipID).Error; err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid Project."}
			}
			project = membership.Project
		} else if taskID != "" {
			var task models.Task
			if err := initializers.DB.Preload("Project").Preload("Project.Memberships").First(&task, "id = ?", taskID).Error; err != nil {
				var subTask models.SubTask
				if err := initializers.DB.Preload("Task").Preload("Task.Project").Preload("Task.Project.Memberships").First(&subTask, "id = ?", taskID).Error; err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Project."}
				}
				project = subTask.Task.Project
			} else {
				project = task.Project
			}
		}

		if project.UserID.String() == loggedInUserID {
			return c.Next()
		}

		var check bool
		for _, membership := range project.Memberships {
			if membership.UserID.String() == loggedInUserID {
				if !checkProjectAccess(membership.Role, Role) {
					return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
				}
				c.Set("projectMemberID", c.GetRespHeader("loggedInUserID"))
				c.Set("loggedInUserID", membership.Project.UserID.String())
				check = true
				break
			}
		}
		if !check {
			return &fiber.Error{Code: 403, Message: "Cannot access this project"}
		}

		return c.Next()
	}
}

func GroupChatAdminAuthorization() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		groupChatID := c.Params("chatID")
		loggedInUserID := c.GetRespHeader("loggedInUserID")

		var chatMembership models.GroupChatMembership
		err := initializers.DB.First(&chatMembership, "group_chat_id = ? AND user_id = ?", groupChatID, loggedInUserID).Error
		if err != nil {
			return &fiber.Error{Code: 400, Message: "No chat of this id found."}
		}

		if chatMembership.Role != models.ChatAdmin {
			return &fiber.Error{Code: 403, Message: "You do not have the permission to perform this action."}
		}

		return c.Next()
	}
}

func TaskUsersCheck(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	taskID := c.Params("taskID")

	var task models.Task
	if err := initializers.DB.Preload("Users").First(&task, "id = ?", taskID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Task of this id found."}
	}

	var check bool
	for _, user := range task.Users {
		if user.ID.String() == loggedInUserID {
			check = true
			break
		}
	}

	if !check {
		return &fiber.Error{Code: 403, Message: "Cannot access this task"}
	}

	return c.Next()
}

func SubTaskUsersAuthorization(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	subTaskID := c.Params("taskID")

	var subTask models.SubTask
	if err := initializers.DB.Preload("Task").Preload("Task.Users").First(&subTask, "id = ?", subTaskID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Sub Task of this id found."}
	}

	var check bool
	for _, user := range subTask.Task.Users {
		if user.ID.String() == loggedInUserID {
			check = true
			break
		}
	}

	if !check {
		return &fiber.Error{Code: 403, Message: "Cannot access this task"}
	}

	return c.Next()
}
