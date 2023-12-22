package project_controllers

import (
	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetNonMembers(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	var project models.Project
	if err := initializers.DB.Where("id = ?", projectID).Preload("Memberships").First(&project).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Project ID"}
	}

	var membershipUserIDs []string

	for _, membership := range project.Memberships {
		membershipUserIDs = append(membershipUserIDs, membership.UserID.String())
	}

	membershipUserIDs = append(membershipUserIDs, project.UserID.String())

	searchedDB := API.Search(c, 0)(initializers.DB)

	var users []models.User
	if err := searchedDB.Where("id NOT IN (?)", membershipUserIDs).Limit(10).Find(&users).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func AddMember(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Project ID"}
	}

	var reqBody struct {
		UserID string
		Title  string
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", reqBody.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No User of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id = ? and user_id=?", parsedProjectID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if reqBody.UserID == project.UserID.String() {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}

	var membership models.Membership
	if err := initializers.DB.Where("user_id=? AND project_id=?", user.ID, parsedProjectID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			var existingInvitation models.Invitation
			err := initializers.DB.Where("user_id=? AND project_id=? AND status=0", user.ID, parsedProjectID).First(&existingInvitation).Error
			if err == nil {
				return &fiber.Error{Code: 400, Message: "Have already invited this User."}
			}

			var invitation models.Invitation
			invitation.ProjectID = &parsedProjectID
			invitation.UserID = user.ID
			invitation.Title = reqBody.Title

			result := initializers.DB.Create(&invitation)

			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			projectMemberID := c.GetRespHeader("projectMemberID")
			parsedID, _ := uuid.Parse(projectMemberID)
			go routines.MarkProjectHistory(project.ID, parsedID, 0, &invitation.UserID, nil, nil, &invitation.ID, nil, "")

			invitation.User = user

			cache.RemoveProject("-workspace--" + project.Slug)

			return c.Status(201).JSON(fiber.Map{
				"status":     "success",
				"message":    "Invitation sent to the user.",
				"invitation": invitation,
			})
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	} else {
		return &fiber.Error{Code: 400, Message: "User is a already a collaborator of this project."}
	}
}

func RemoveMember(c *fiber.Ctx) error { //TODO add manager cannot remove manager
	membershipID := c.Params("membershipID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var membership models.Membership
	if err := initializers.DB.Preload("Project").First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if membership.Project.UserID != parsedLoggedInUserID {
		return &fiber.Error{Code: 403, Message: "You do not have the permission to perform this action."}
	}

	err = processLeaveProject(&membership)
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	parsedUserID := membership.UserID
	parsedProjectID := membership.ProjectID
	projectSlug := membership.Project.Slug

	projectMemberID := c.GetRespHeader("projectMemberID")
	parsedID, _ := uuid.Parse(projectMemberID)
	go routines.MarkProjectHistory(parsedProjectID, parsedID, 11, &parsedUserID, nil, nil, nil, nil, membership.Title)

	cache.RemoveProject(projectSlug)
	cache.RemoveProject("-workspace--" + projectSlug)

	go routines.DecrementProjectMember(parsedProjectID)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "User removed to the project.",
	})
}

func LeaveProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var membership models.Membership
	if err := initializers.DB.Preload("Project").First(&membership, "project_id = ? && user_id=?", projectID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	err := processLeaveProject(&membership)
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	parsedUserID := membership.UserID
	parsedProjectID := membership.ProjectID
	projectSlug := membership.Project.Slug

	go routines.MarkProjectHistory(parsedProjectID, parsedUserID, 10, nil, nil, nil, nil, nil, membership.Title)

	cache.RemoveProject(projectSlug)
	cache.RemoveProject("-workspace--" + projectSlug)

	go routines.DecrementProjectMember(parsedProjectID)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "You left the project.",
	})
}

func ChangeMemberRole(c *fiber.Ctx) error {
	membershipID := c.Params("membershipID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	parsedMembershipID, err := uuid.Parse(membershipID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Membership ID"}
	}

	var reqBody struct {
		Title string             `json:"title"`
		Role  models.ProjectRole `json:"role"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	var membership models.Membership
	if err := initializers.DB.Preload("Project").First(&membership, "id = ?", parsedMembershipID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Membership of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if parsedLoggedInUserID != membership.Project.UserID {
		var updatingUserMembership models.Membership
		if err := initializers.DB.Preload("Project").First(&updatingUserMembership, "project_id = ? AND user_id = ?", membership.ProjectID, parsedLoggedInUserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &fiber.Error{Code: 403, Message: "You are not a part of this project."}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}

		if updatingUserMembership.Role != models.ProjectManager {
			return &fiber.Error{Code: 403, Message: "Cannot perform this action."}
		}
		if membership.Role == models.ProjectManager {
			return &fiber.Error{Code: 403, Message: "Cannot perform this action."}
		}
		if reqBody.Role == models.ProjectManager {
			return &fiber.Error{Code: 403, Message: "Cannot perform this action."}
		}
	}

	membership.Role = reqBody.Role

	if reqBody.Title != "" {
		membership.Title = reqBody.Title
	}

	result := initializers.DB.Save(&membership)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	cache.RemoveProject(membership.Project.Slug)
	cache.RemoveProject("-workspace--" + membership.Project.Slug)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User membership updated.",
	})
}

func processLeaveProject(membership *models.Membership) error {
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

	// Step 1: Retrieve the user's group chat memberships in the specified project
	var memberships []models.GroupChatMembership
	if err := tx.Where("user_id = ? AND group_chat_id IN (SELECT id FROM group_chats WHERE project_id = ?)", membership.UserID, membership.ProjectID).Find(&memberships).Error; err != nil {
		return err
	}

	// Step 2: Delete the group chat memberships
	if err := tx.Delete(&memberships).Error; err != nil {
		return err
	}

	// Step 3: Find all tasks assigned to the user in the given project
	var tasks []models.Task
	if err := tx.
		Joins("JOIN task_assigned_users ON tasks.id = task_assigned_users.task_id").
		Where("tasks.project_id = ? AND task_assigned_users.user_id = ?", membership.ProjectID, membership.UserID).
		Find(&tasks).Error; err != nil {
		return err
	}

	// Step 4: Remove the user from the assigned users of each task
	for _, task := range tasks {
		if err := tx.Model(&task).Association("Users").Delete(&models.User{ID: membership.UserID}); err != nil {
			return err
		}
	}

	// Step 5: Find all subtasks assigned to the user in the given project
	var subtasks []models.SubTask
	if err := tx.
		Joins("JOIN tasks ON sub_tasks.task_id = tasks.id").
		Joins("JOIN task_assigned_users ON tasks.id = task_assigned_users.task_id").
		Where("tasks.project_id = ? AND task_assigned_users.user_id = ?", membership.ProjectID, membership.UserID).
		Find(&subtasks).Error; err != nil {
		return err
	}

	// Step 6: Remove the user from the assigned users of each subtask
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
