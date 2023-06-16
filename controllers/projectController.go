package controllers

import (
	"reflect"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var project models.Project
	if err := initializers.DB.Preload("User").First(&project, "id = ?", parsedProjectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"project": project,
	})
}

func GetUserProjects(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var projects []models.Project
	if err := initializers.DB.Where("user_id = ? AND is_private = ?", userID, false).Find(&projects).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"projects": projects,
	})
}

func GetUserContributingProjects(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var memberships []models.Membership
	if err := initializers.DB.Preload("Project").Select("project_id").Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var projects []models.Project
	for _, membership := range memberships {
		projects = append(projects, membership.Project)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"projects": projects,
	})
}

func GetProjectContributors(c *fiber.Ctx) error { //! Add search here
	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var memberships []models.Membership
	if err := initializers.DB.Preload("User").Where("project_id = ?", parsedProjectID).Find(&memberships).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var users []models.User
	for _, membership := range memberships {
		users = append(users, membership.User)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"users":   users,
	})
}

func AddProject(c *fiber.Ctx) error {
	var reqBody schemas.ProjectCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.ProjectCreateSchema](reqBody); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(c.GetRespHeader("loggedInUserID"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	picName, err := utils.SaveFile(c, "coverPic", "projects/coverPics", true, 900, 400)
	if err != nil {
		return err
	}

	newProject := models.Project{
		UserID:      parsedID,
		Title:       reqBody.Title,
		Tagline:     reqBody.Tagline,
		CoverPic:    picName,
		Description: reqBody.Description,
		Tags:        reqBody.Tags,
		Category:    reqBody.Category,
		IsPrivate:   reqBody.IsPrivate,
		Links:       reqBody.Links,
	}

	result := initializers.DB.Create(&newProject)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating project"}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Project Added",
		"project": newProject,
	})
}

func UpdateProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", parsedProjectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	var reqBody schemas.ProjectUpdateSchema
	c.BodyParser(&reqBody)

	picName, err := utils.SaveFile(c, "coverPic", "projects/coverPics", true, 900, 400)
	if err != nil {
		return err
	}
	reqBody.CoverPic = picName

	// if reqBody.Tagline != "" {
	// 	project.Tagline = reqBody.Tagline
	// }
	// if reqBody.CoverPic != "" {
	// 	project.CoverPic = reqBody.CoverPic
	// }
	// if reqBody.Description != "" {
	// 	project.Description = reqBody.Description
	// }
	// if reqBody.Page != "" {
	// 	project.Page = reqBody.Page
	// }
	// if len(reqBody.Tags) != 0 {
	// 	project.Tags = reqBody.Tags
	// }
	// if reqBody.IsPrivate {
	// 	project.IsPrivate = true
	// } else {
	// 	project.IsPrivate = false
	// }
	// if len(reqBody.Links) != 0 {
	// 	project.Links = reqBody.Links
	// }
	// if len(reqBody.PrivateLinks) != 0 {
	// 	project.PrivateLinks = reqBody.PrivateLinks
	// }

	projectValue := reflect.ValueOf(project).Elem()
	reqBodyValue := reflect.ValueOf(reqBody)

	for i := 0; i < reqBodyValue.NumField(); i++ {
		field := reqBodyValue.Type().Field(i)
		fieldValue := reqBodyValue.Field(i)

		if fieldValue.IsZero() {
			continue
		}

		projectField := projectValue.FieldByName(field.Name)
		if projectField.IsValid() && projectField.CanSet() {
			projectField.Set(fieldValue)
		}
	}

	if err := initializers.DB.Save(&project).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Project updated successfully",
		"project": project,
	})
}

func DeleteProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", parsedProjectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	if err := initializers.DB.Delete(&project).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Project deleted successfully",
	})
}
