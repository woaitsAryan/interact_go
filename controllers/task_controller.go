package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetTask(taskType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		taskID := c.Params("taskID")

		parsedTaskID, err := uuid.Parse(taskID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID"}
		}

		switch taskType {
		case "task":
			var task models.Task
			if err := initializers.DB.
				Preload("User").
				Preload("SubTask").
				First(&task, "id = ?", parsedTaskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			return c.Status(200).JSON(fiber.Map{
				"status":  "success",
				"message": "",
				"task":    task,
			})
		case "subtask":
			var task models.SubTask
			if err := initializers.DB.
				Preload("User").
				First(&task, "id = ?", parsedTaskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Sub Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			return c.Status(200).JSON(fiber.Map{
				"status":  "success",
				"message": "",
				"task":    task,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": config.SERVER_ERROR,
		})
	}
}

func AddTask(taskType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqBody schemas.TaskCreateSchema
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: err.Error()}
		}

		if err := helpers.Validate[schemas.TaskCreateSchema](reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: err.Error()}
		}

		var users []models.User

		switch taskType {
		case "task":
			projectID := c.Params("projectID")

			var project models.Project
			if err := initializers.DB.Preload("Memberships").Preload("Memberships.User").First(&project, "id = ?", projectID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			var projectMembers []models.User
			for _, membership := range project.Memberships {
				projectMembers = append(projectMembers, membership.User)
			}

			for _, userID := range reqBody.Users {
				if GetUserIndex(userID, projectMembers) != -1 || project.UserID.String() == userID {
					var user models.User
					if err := initializers.DB.First(&user, "id = ?", userID).Error; err == nil {
						users = append(users, user)
					}
				}
			}

			task := models.Task{
				ProjectID:   project.ID,
				Title:       reqBody.Title,
				Description: reqBody.Description,
				Tags:        reqBody.Tags,
				Deadline:    reqBody.Dateline,
				Users:       users,
			}

			result := initializers.DB.Create(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			return c.Status(201).JSON(fiber.Map{
				"status":  "success",
				"message": "",
				"task":    task,
			})

		case "subtask":
			taskID := c.Params("taskID")

			var task models.Task
			if err := initializers.DB.Preload("Users").First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			for _, userID := range reqBody.Users {
				if GetUserIndex(userID, task.Users) != -1 {
					var user models.User
					if err := initializers.DB.First(&user, "id = ?", userID).Error; err == nil {
						users = append(users, user)
					}
				}
			}

			subTask := models.SubTask{
				TaskID:      task.ID,
				Title:       reqBody.Title,
				Description: reqBody.Description,
				Tags:        reqBody.Tags,
				Deadline:    reqBody.Dateline,
				Users:       users,
			}

			result := initializers.DB.Create(&subTask)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			return c.Status(201).JSON(fiber.Map{
				"status":  "success",
				"message": "",
				"task":    subTask,
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": config.SERVER_ERROR,
		})
	}
}

func EditTask(taskType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqBody schemas.TaskEditSchema
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: err.Error()}
		}

		if err := helpers.Validate[schemas.TaskEditSchema](reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: err.Error()}
		}

		taskID := c.Params("taskID")

		switch taskType {
		case "task":
			var task models.Task
			if err := initializers.DB.First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			if reqBody.Title != "" {
				task.Title = reqBody.Title
			}
			if reqBody.Description != "" {
				task.Description = reqBody.Description
			}
			if reqBody.Tags != nil {
				task.Tags = reqBody.Tags
			}

			result := initializers.DB.Save(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

		case "subtask":
			var task models.SubTask
			if err := initializers.DB.First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Sub Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			if reqBody.Title != "" {
				task.Title = reqBody.Title
			}
			if reqBody.Description != "" {
				task.Description = reqBody.Description
			}
			if reqBody.Tags != nil {
				task.Tags = reqBody.Tags
			}

			result := initializers.DB.Save(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Task edited",
		})
	}
}

func MarkTaskCompleted(taskType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqBody struct {
			IsCompleted bool `json:"isCompleted"`
		}
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		taskID := c.Params("taskID")
		userID := c.GetRespHeader("loggedInUserID")

		switch taskType {
		case "task":
			var task models.Task
			if err := initializers.DB.Preload("Users").First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			userIndex := GetUserIndex(userID, task.Users)
			if userIndex == -1 {
				return &fiber.Error{Code: 403, Message: "Cannot Perform this action"}
			}

			task.IsCompleted = reqBody.IsCompleted

			result := initializers.DB.Save(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			// if reqBody.IsCompleted{
			// 	go MarkSubTasksCompleted(task.ID)
			// }

		case "subtask":
			var task models.SubTask
			if err := initializers.DB.Preload("Users").First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Sub Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			userIndex := GetUserIndex(userID, task.Users)
			if userIndex == -1 {
				return &fiber.Error{Code: 403, Message: "Cannot Perform this action"}
			}

			task.IsCompleted = reqBody.IsCompleted

			result := initializers.DB.Save(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Task edited",
		})
	}
}

func DeleteTask(taskType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		taskID := c.Params("taskID")

		switch taskType {
		case "task":
			var task models.Task
			if err := initializers.DB.
				First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			// Delete all users from the task_assigned_users table
			if err := initializers.DB.Model(&task).Association("Users").Clear(); err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			result := initializers.DB.Delete(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

		case "subtask":
			var task models.SubTask
			if err := initializers.DB.
				First(&task, "id = ?", taskID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return &fiber.Error{Code: 400, Message: "No Sub Task of this ID found."}
				}
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			// Delete all users from the subtask_assigned_users table
			if err := initializers.DB.Model(&task).Association("Users").Clear(); err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			result := initializers.DB.Delete(&task)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		}

		return c.Status(204).JSON(fiber.Map{
			"status":  "success",
			"message": "Task Deleted",
		})
	}
}
