package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddTaskUser(taskType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		taskID := c.Params("taskID")

		var reqBody struct {
			UserID string `json:"userID"`
		}
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		var user models.User
		if err := initializers.DB.First(&user, "id = ?", reqBody.UserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &fiber.Error{Code: 400, Message: "No User of this ID found."}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		switch taskType {
		case "task":
			var task models.Task
			if err := initializers.DB.First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			var project models.Project
			if err := initializers.DB.Preload("Memberships").Preload("Memberships.User").First(&project, "id = ?", task.ProjectID).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			check := false
			for _, membership := range project.Memberships {
				if membership.UserID == user.ID {
					check = true
					break
				}
			}

			if check || project.UserID == user.ID {
				task.Users = append(task.Users, user)
			} else {
				return &fiber.Error{Code: 400, Message: "User not a member of this Project."}
			}

			result := initializers.DB.Save(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
			}
		case "org_task":
			var task models.Task
			if err := initializers.DB.First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			var organization models.Organization
			if err := initializers.DB.Preload("Memberships").Preload("Memberships.User").First(&organization, "id = ?", task.OrganizationID).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			check := false
			for _, membership := range organization.Memberships {
				if membership.UserID == user.ID {
					check = true
					break
				}
			}

			if check {
				task.Users = append(task.Users, user)
			} else {
				return &fiber.Error{Code: 400, Message: "User not a member of this Project."}
			}

			result := initializers.DB.Save(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
			}
		case "subtask":
			var subTask models.SubTask
			if err := initializers.DB.First(&subTask, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Sub Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			var task models.Task
			if err := initializers.DB.Preload("Users").First(&task, "id = ?", subTask.TaskID).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			check := false
			for _, taskUser := range task.Users {
				if taskUser.ID == user.ID {
					check = true
					break
				}
			}

			if check {
				subTask.Users = append(subTask.Users, user)
			} else {
				return &fiber.Error{Code: 400, Message: "User not a member of this Task."}
			}

			result := initializers.DB.Save(&subTask)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
			}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User added",
		})
	}
}

func GetUserIndex(userID string, users []models.User) int {
	var userIndex = -1
	for i, u := range users {
		if u.ID.String() == userID {
			userIndex = i
			break
		}
	}

	return userIndex
}

func RemoveTaskUser(taskType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		taskID := c.Params("taskID")
		userID := c.Params("userID")

		var user models.User
		if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &fiber.Error{Code: 400, Message: "No User of this ID found."}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		switch taskType {
		case "task":
			var task models.Task
			if err := initializers.DB.Preload("Users").First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			userIndex := GetUserIndex(userID, task.Users)
			if userIndex == -1 {
				return &fiber.Error{Code: 400, Message: "User not assigned to this task."}
			}

			task.Users = append(task.Users[:userIndex], task.Users[userIndex+1:]...)
			result := initializers.DB.Save(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
			}

			// Delete the user from the task_assigned_users table
			if err := initializers.DB.Model(&task).Association("Users").Delete(&user); err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}
		case "subtask":
			var task models.SubTask
			if err := initializers.DB.Preload("Users").First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Sub Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			userIndex := GetUserIndex(userID, task.Users)
			if userIndex == -1 {
				return &fiber.Error{Code: 400, Message: "User not assigned to this task."}
			}

			task.Users = append(task.Users[:userIndex], task.Users[userIndex+1:]...)

			result := initializers.DB.Save(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
			}

			// Delete the user from the subtask_assigned_users table
			if err := initializers.DB.Model(&task).Association("Users").Delete(&user); err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User removed",
		})
	}
}
