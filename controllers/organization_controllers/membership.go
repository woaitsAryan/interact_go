package organization_controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/controllers/auth_controllers"
	"github.com/Pratham-Mishra04/interact/controllers/project_controllers"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/Pratham-Mishra04/interact/utils/select_fields"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetExploreMemberships(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var memberships []models.OrganizationMembership
	if err := paginatedDB.Where("organization_id = ?", orgID).Preload("User").Find(&memberships).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID"}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"memberships": memberships,
	})
}

func GetNonMembers(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var organization models.Organization
	if err := initializers.DB.Where("id = ?", orgID).Preload("Memberships").First(&organization).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID"}
	}

	var membershipUserIDs []string

	for _, membership := range organization.Memberships {
		membershipUserIDs = append(membershipUserIDs, membership.UserID.String())
	}

	membershipUserIDs = append(membershipUserIDs, organization.UserID.String())

	searchedDB := API.Search(c, 0)(initializers.DB)

	var users []models.User
	if err := searchedDB.Where("id NOT IN (?)", membershipUserIDs).
		Where("active=? AND onboarding_completed=?", true, true).
		Where("verified=?", true).
		Where("username != email").
		Where("organization_status=?", false).
		Where("username != users.email").
		Limit(10).
		Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetMemberships(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var organization models.Organization
	if err := initializers.DB.Where("id = ?", orgID).
		Preload("Memberships", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Memberships.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.ExtendedUser)
		}).
		Preload("Invitations", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Invitations.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.ExtendedUser)
		}).
		First(&organization).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID"}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"organization": organization,
	})
}

func AddMember(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	orgID := c.Params("orgID")

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}

	parsedOrganizationID, err := uuid.Parse(orgID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID"}
	}

	var reqBody struct {
		UserID string `json:"userID"`
		Title  string `json:"title"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ? AND organization_status=false", reqBody.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No User of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var organization models.Organization
	if err := initializers.DB.First(&organization, "id = ? and user_id=?", parsedOrganizationID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Organization of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.UserID == organization.UserID.String() {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}

	var membership models.OrganizationMembership
	if err := initializers.DB.Where("user_id=? AND organization_id=?", user.ID, parsedOrganizationID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var existingInvitation models.Invitation
			err := initializers.DB.Where("user_id=? AND organization_id=?", user.ID, parsedOrganizationID).First(&existingInvitation).Error
			if err == nil {
				if existingInvitation.Status != 0 {
					existingInvitation.Status = 0
					existingInvitation.Title = reqBody.Title
					existingInvitation.CreatedAt = time.Now()

					if err := initializers.DB.Save(&existingInvitation).Error; err != nil {
						return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
					}

					existingInvitation.User = user

					return c.Status(201).JSON(fiber.Map{
						"status":     "success",
						"message":    "Invitation sent to the user.",
						"invitation": existingInvitation,
					})
				}
				return &fiber.Error{Code: 400, Message: "An invitation is already sent to this user."}
			}

			var invitation models.Invitation
			invitation.OrganizationID = &parsedOrganizationID
			invitation.UserID = user.ID
			invitation.Title = reqBody.Title

			result := initializers.DB.Create(&invitation)

			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			invitation.User = user

			go routines.MarkOrganizationHistory(parsedOrganizationID, parsedUserID, 3, nil, nil, nil, nil, &invitation.ID, nil, nil, nil, nil, nil, "")
			go cache.RemoveOrganization("-access--" + organization.ID.String())
			return c.Status(201).JSON(fiber.Map{
				"status":     "success",
				"message":    "Invitation sent to the user.",
				"invitation": invitation,
			})
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	} else {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}
}

func RemoveMember(c *fiber.Ctx) error {
	membershipID := c.Params("membershipID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	orgMemberID := c.GetRespHeader("orgMemberID")
	parsedOrgMemberID, _ := uuid.Parse(orgMemberID)
	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var membership models.OrganizationMembership
	if err := initializers.DB.Preload("Organization").First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if membership.UserID.String() == orgMemberID {
		return &fiber.Error{Code: 400, Message: "Cannot remove yourself using this route."}
	}

	if membership.Organization.UserID != parsedLoggedInUserID {
		return &fiber.Error{Code: 403, Message: "You do not have the permission to perform this action."}
	}

	err = processLeaveOrganization(&membership)
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.DecrementOrgMember(membership.OrganizationID)
	go routines.MarkOrganizationHistory(membership.OrganizationID, parsedOrgMemberID, 5, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, membership.Title)
	go cache.RemoveOrganization("-access--" + membership.OrganizationID.String())

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "User removed to the project.",
	})
}

func LeaveOrganization(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	orgMemberID := c.GetRespHeader("orgMemberID")
	parsedOrgMemberID, _ := uuid.Parse(orgMemberID)

	var membership models.OrganizationMembership
	if err := initializers.DB.Preload("Organization").First(&membership, "user_id=? AND organization_id = ?", orgMemberID, orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var reqBody struct {
		VerificationCode string `json:"otp"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "OTP not provided."}
	}

	data, err := cache.GetOtpFromCache(membership.UserID.String() + "-" + membership.ID.String())
	if err != nil {
		return &fiber.Error{Code: 400, Message: "OTP Expired"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data), []byte(reqBody.VerificationCode)); err != nil {
		return &fiber.Error{Code: 400, Message: "Incorrect OTP"}
	}

	err = processLeaveOrganization(&membership)
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.DecrementOrgMember(membership.OrganizationID)
	go routines.MarkOrganizationHistory(membership.OrganizationID, parsedOrgMemberID, 15, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, membership.Title)
	go cache.RemoveOrganization("-access--" + membership.OrganizationID.String())

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "You left the organization",
	})
}

func ChangeMemberRole(c *fiber.Ctx) error {
	membershipID := c.Params("membershipID")

	orgChangedUserID := c.GetRespHeader("loggedInUserID")
	loggedInUserID := c.GetRespHeader("orgMemberID")

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var reqBody struct {
		Role  models.OrganizationRole `json:"role"`
		Title string                  `json:"title"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var membership models.OrganizationMembership
	if err := initializers.DB.First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	membership.Title = reqBody.Title

	if orgChangedUserID == loggedInUserID {
		membership.Role = reqBody.Role
	} else {
		if reqBody.Role != models.Manager && membership.Role != models.Manager {
			membership.Role = reqBody.Role
		} else {
			return &fiber.Error{Code: 403, Message: "You don't have the privileges to perform this action."}
		}
	}

	result := initializers.DB.Save(&membership)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	parsedUserID, _ := uuid.Parse(loggedInUserID)
	go routines.MarkOrganizationHistory(membership.OrganizationID, parsedUserID, 30, nil, nil, nil, nil, nil, nil, nil, nil, nil, &membership.ID, "")

	go cache.RemoveOrganization("-access--" + membership.OrganizationID.String())

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User membership updated.",
	})
}

func processLeaveOrganization(membership *models.OrganizationMembership) error {
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback() // Rollback the transaction on panic
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "completeLeaveProject")
		}
	}()

	// Step 1: Retrieve the user's group chat memberships in the specified org
	var chatMemberships []models.GroupChatMembership
	if err := tx.Where("user_id = ? AND group_chat_id IN (SELECT id FROM group_chats WHERE organization_id = ?)", membership.UserID, membership.OrganizationID).Find(&chatMemberships).Error; err != nil {
		return err
	}

	// Step 2: Delete the group chat memberships
	for _, chatMembership := range chatMemberships {
		if err := tx.Delete(&chatMembership).Error; err != nil {
			return err
		}
	}

	// Step 3: Retrieve the user's project memberships in the specified org
	var projectMemberships []models.Membership
	if err := tx.Where("user_id = ? AND project_id IN (SELECT id FROM projects WHERE user_id = ?)", membership.UserID, membership.Organization.UserID).Find(&projectMemberships).Error; err != nil {
		return err
	}

	// Step 4: Delete project memberships
	for _, membership := range projectMemberships {
		project_controllers.ProcessLeaveProject(&membership)
	}

	// Step 5: Find all tasks assigned to the user in the given org
	var tasks []models.Task
	if err := tx.
		Joins("JOIN task_assigned_users ON tasks.id = task_assigned_users.task_id").
		Where("tasks.organization_id = ? AND task_assigned_users.user_id = ?", membership.OrganizationID, membership.UserID).
		Find(&tasks).Error; err != nil {
		return err
	}

	// Step 6: Remove the user from the assigned users of each task
	for _, task := range tasks {
		if err := tx.Model(&task).Association("Users").Delete(&models.User{ID: membership.UserID}); err != nil {
			return err
		}
	}

	// Step 7: Find all subtasks assigned to the user in the given org
	var subtasks []models.SubTask
	if err := tx.
		Joins("JOIN tasks ON sub_tasks.task_id = tasks.id").
		Joins("JOIN sub_task_assigned_users ON tasks.id = sub_task_assigned_users.sub_task_id").
		Where("tasks.organization_id = ? AND sub_task_assigned_users.user_id = ?", membership.OrganizationID, membership.UserID).
		Find(&subtasks).Error; err != nil {
		return err
	}

	// Step 8: Remove the user from the assigned users of each subtask
	for _, subtask := range subtasks {
		if err := tx.Model(&subtask).Association("Users").Delete(&models.User{ID: membership.UserID}); err != nil {
			return err
		}
	}

	result := tx.Delete(&membership)
	if result.Error != nil {
		return result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func SendLeaveOrgVerificationCode(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("orgMemberID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	orgID := c.Params("orgID")

	code := auth_controllers.GenerateOTP(6)
	hash, err := bcrypt.GenerateFromPassword([]byte(code), 10)
	if err != nil {
		go helpers.LogServerError("Error while hashing an OTP.", err, c.Path())
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	var membership models.OrganizationMembership
	if err := initializers.DB.Preload("User").Where("organization_id=? AND user_id=?", orgID, parsedLoggedInUserID).First(&membership).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Membership found."}
	}

	err = helpers.SendMailReq(membership.User.Email, config.OTP_VERIFICATION_MAIL, &membership.User, &code, nil)
	if err != nil {
		return &fiber.Error{Code: 500, Message: config.SERVER_ERROR}
	}

	err = cache.SetOtpToCache(membership.UserID.String()+"-"+membership.ID.String(), []byte(hash))
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "OTP sent to registered mail",
	})
}
