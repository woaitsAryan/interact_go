package middlewares

import (
	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func checkOrgAccess(UserRole models.OrganizationRole, AuthorizedRole models.OrganizationRole) bool {
	if UserRole == models.Manager {
		return true
	} else if UserRole == models.Senior {
		return AuthorizedRole != models.Manager
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

		var organization models.Organization
		if orgID != "" {
			orgInCache, err := cache.GetOrganization("-access--" + orgID)
			if err == nil {
				organization = *orgInCache
			} else {
				if err := initializers.DB.Preload("Memberships").First(&organization, "id=?", orgID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return &fiber.Error{Code: 401, Message: "No Organization of this ID Found."}
					}
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}

				go cache.SetOrganization("-access--"+organization.ID.String(), &organization)
			}
		} else {
			return &fiber.Error{Code: 401, Message: "Invalid Organization ID."}
		}

		if organization.UserID.String() == loggedInUserID {
			c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
			return c.Next()
		}

		var check bool
		check = false

		for _, membership := range organization.Memberships {
			if membership.UserID.String() == loggedInUserID {
				if !checkOrgAccess(membership.Role, Role) {
					return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
				}
				c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
				c.Set("loggedInUserID", organization.UserID.String())
				check = true
				break
			}
		}
		if !check {
			return &fiber.Error{Code: 403, Message: "Cannot access this Organization."}
		}

		return c.Next()
	}
}

func OrgEventRoleAuthorization(Role models.OrganizationRole) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		orgID := c.Params("orgID")
		eventID := c.Params("eventID")

		var event models.Event
		if err := initializers.DB.Preload("Organization").
				  Preload("Organization.Memberships").
				  Preload("CoOwnedBy").
				  Preload("CoOwnedBy.Memberships").
				  First(&event, "id = ?", eventID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No Event of this id found."}
		}

		var organization models.Organization
		if orgID != "" {
			orgInCache, err := cache.GetOrganization("-access--" + orgID)
			if err == nil {
				organization = *orgInCache
			} else {
				if err := initializers.DB.Preload("Memberships").First(&organization, "id=?", orgID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return &fiber.Error{Code: 401, Message: "No Organization of this ID Found."}
					}
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}

				go cache.SetOrganization("-access--"+organization.ID.String(), &organization)
			}
		} else {
			return &fiber.Error{Code: 401, Message: "Invalid Organization ID."}
		}

		if organization.UserID.String() == loggedInUserID {
			c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
			return c.Next()
		}

		var check bool
		check = false

		for _, membership := range organization.Memberships {
			if membership.UserID.String() == loggedInUserID {
				if !checkOrgAccess(membership.Role, Role) {
					return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
				}
				c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
				c.Set("loggedInUserID", organization.UserID.String())
				check = true
				break
			}
		}

		for _, coOwnOrganization := range event.CoOwnedBy {
			for _, membership := range coOwnOrganization.Memberships {
				if membership.UserID.String() == loggedInUserID {
					if !checkOrgAccess(membership.Role, Role) {
						return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
					}
					c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
					c.Set("loggedInUserID", organization.UserID.String())
					check = true
					break
				}
			}
		}

		if !check {
			return &fiber.Error{Code: 403, Message: "Cannot access this Organization."}
		}

		return c.Next()
	}
}

func OrgPollAuthorization() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		orgID := c.Params("orgID")
		pollID := c.Params("pollID")

		var organization models.Organization
		if orgID != "" {
			orgInCache, err := cache.GetOrganization("-access--" + orgID)
			if err == nil {
				organization = *orgInCache
			} else {
				if err := initializers.DB.Preload("Memberships").First(&organization, "id=?", orgID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return &fiber.Error{Code: 401, Message: "No Organization of this ID Found."}
					}
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}

				go cache.SetOrganization("-access--"+organization.ID.String(), &organization)
			}
		} else {
			return &fiber.Error{Code: 401, Message: "Invalid Organization ID."}
		}
		var poll models.Poll
		if err := initializers.DB.First(&poll, "id = ?", pollID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No poll of this id found."}
		}
		if poll.IsOpen {
			c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
			return c.Next()
		}

		if organization.UserID.String() == loggedInUserID {
			c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
			return c.Next()
		}

		var check bool
		check = false

		for _, membership := range organization.Memberships {
			if membership.UserID.String() == loggedInUserID {
				if !checkOrgAccess(membership.Role, models.Member) {
					return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
				}
				c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
				c.Set("loggedInUserID", organization.UserID.String())
				check = true
				break
			}
		}
		if !check {
			return &fiber.Error{Code: 403, Message: "Cannot access this Organization."}
		}
		return c.Next()
	}
}

func OrgBucketAuthorization(action string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		orgID := c.Params("orgID")
		resourceBucketID := c.Params("resourceBucketID")

		var organization models.Organization
		if orgID != "" {
			orgInCache, err := cache.GetOrganization("-access--" + orgID)
			if err == nil {
				organization = *orgInCache
			} else {
				if err := initializers.DB.Preload("Memberships").First(&organization, "id=?", orgID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return &fiber.Error{Code: 401, Message: "No Organization of this ID Found."}
					}
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}

				go cache.SetOrganization("-access--"+organization.ID.String(), &organization)
			}
		} else {
			return &fiber.Error{Code: 401, Message: "Invalid Organization ID."}
		}

		if organization.UserID.String() == loggedInUserID {
			c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
			return c.Next()
		}

		var resourceBucket models.ResourceBucket

		resourceBucketInCache, err := cache.GetResourceBucket(resourceBucketID)
		if err == nil {
			resourceBucket = *resourceBucketInCache
		} else {
			if err := initializers.DB.Where("id=? AND organization_id = ?", resourceBucketID, orgID).First(&resourceBucket).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource Bucket does not exist."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}
		}

		var check bool
		check = false

		level := models.Member
		if action == "view" {
			level = resourceBucket.ViewAccess
		} else {
			level = resourceBucket.EditAccess
		}

		for _, membership := range organization.Memberships {
			if membership.UserID.String() == loggedInUserID {
				if !checkOrgAccess(membership.Role, level) {
					return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
				}
				c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
				c.Set("loggedInUserID", organization.UserID.String())
				check = true
				break
			}
		}
		if !check {
			return &fiber.Error{Code: 403, Message: "Cannot access this Organization."}
		}

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

		projectInCache, err := cache.GetProject("-access--" + projectID)
		if err == nil {
			project = *projectInCache
		} else {
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

			go cache.SetProject("-access--"+project.ID.String(), &project)
		}

		if project.UserID.String() == loggedInUserID {
			c.Set("projectMemberID", c.GetRespHeader("loggedInUserID"))
			return c.Next()
		}

		var check bool
		check = false

		for _, membership := range project.Memberships {
			if membership.UserID.String() == loggedInUserID {
				if !checkProjectAccess(membership.Role, Role) {
					return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
				}
				c.Set("projectMemberID", c.GetRespHeader("loggedInUserID"))
				c.Set("loggedInUserID", project.UserID.String())
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
		err := initializers.DB.Preload("GroupChat.Project").
			Preload("GroupChat.Project.Memberships").
			Preload("GroupChat.Organization.Memberships").
			First(&chatMembership, "group_chat_id = ? AND user_id = ?", groupChatID, loggedInUserID).Error
		if err != nil {
			return &fiber.Error{Code: 400, Message: "No chat of this id found."}
		}

		var accessGranted = false
		if chatMembership.Role == models.ChatAdmin {
			accessGranted = true
		} else if chatMembership.GroupChat.ProjectID != nil {
			roleMemberships := chatMembership.GroupChat.Project.Memberships
			for _, membership := range roleMemberships {
				if membership.UserID.String() == loggedInUserID &&
					(membership.Role == models.ProjectEditor || membership.Role == models.ProjectManager) {
					accessGranted = true
					break
				}
			}
		} else if chatMembership.GroupChat.OrganizationID != nil {
			organizationMemberships := chatMembership.GroupChat.Organization.Memberships
			for _, membership := range organizationMemberships {
				if membership.UserID.String() == loggedInUserID &&
					(membership.Role == models.Manager) {
					accessGranted = true
					break
				}
			}
		}

		if !accessGranted {
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
