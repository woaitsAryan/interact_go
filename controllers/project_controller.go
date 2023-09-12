package controllers

import (
	"log"
	"reflect"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func GetProject(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var project models.Project
	if err := initializers.DB.Omit("private_links").Preload("User").Preload("Openings").Preload("Memberships").First(&project, "slug = ? AND is_private = ? ", slug, false).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var memberships []models.Membership
	if err := initializers.DB.Preload("User").Find(&memberships, "project_id = ?", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	go routines.UpdateProjectViews(&project)

	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, err := uuid.Parse(loggedInUserID)

	if err == nil && parsedLoggedInUserID != project.UserID {
		go routines.UpdateLastViewed(parsedLoggedInUserID, project.ID)
	}

	_, count, err := utils.GetProjectViews(project.ID)
	if err != nil {
		return err
	}
	project.Views = count
	project.Memberships = memberships

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"project": project,
	})
}

func GetWorkSpaceProject(c *fiber.Ctx) error {
	slug := c.Params("slug")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedLoggedInUserID, err := uuid.Parse(loggedInUserID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var project models.Project
	if err := initializers.DB.Preload("User").First(&project, "slug = ?", slug).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var memberships []models.Membership
	if err := initializers.DB.Preload("User").Find(&memberships, "project_id = ?", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var membershipCheck bool
	for _, membership := range memberships {
		if (membership.UserID) == parsedLoggedInUserID {
			membershipCheck = true
		}
	}
	if !membershipCheck && project.UserID != parsedLoggedInUserID {
		return &fiber.Error{Code: 403, Message: "Cannot perform this action."}
	}

	var invitations []models.Invitation
	if err := initializers.DB.Preload("User").Find(&invitations, "project_id = ? AND (status = 0 OR status = -1)", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var chats []models.GroupChat
	if err := initializers.DB.Find(&chats, "project_id = ? ", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	_, count, err := utils.GetProjectViews(project.ID)
	if err != nil {
		return err
	}
	project.Views = count
	project.Memberships = memberships
	project.Invitations = invitations
	project.Chats = chats

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"project": project,
	})
}

func GetMyLikedProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var projectLikes []models.Like
	if err := initializers.DB.Where("user_id = ? AND project_id IS NOT NULL", loggedInUserID).Find(&projectLikes).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var projectIDs []string
	for _, projectLike := range projectLikes {
		projectIDs = append(projectIDs, projectLike.ProjectID.String())
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"projects": projectIDs,
	})
}

func GetUserProjects(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var projects []models.Project
	if err := initializers.DB.Where("user_id = ? AND is_private = ?", userID, false).Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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

func GetProjectContributors(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var memberships []models.Membership
	if err := initializers.DB.Preload("User").Where("project_id = ?", parsedProjectID).Find(&memberships).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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

	slug := slug.Make(reqBody.Title)

	var existingProject models.Project
	if err := initializers.DB.Where("slug=?", slug).First(&existingProject).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			parsedID, err := uuid.Parse(c.GetRespHeader("loggedInUserID"))
			if err != nil {
				return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
			}

			var user models.User
			if err := initializers.DB.Where("id=?", parsedID).First(&user).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}
			if !user.Verified && initializers.CONFIG.ENV == "production" {
				return &fiber.Error{Code: 401, Message: config.VERIFICATION_ERROR}
			}

			picName, err := utils.SaveFile(c, "coverPic", "project/coverPics", true, 1080, 1080)
			if err != nil {
				return err
			}

			newProject := models.Project{
				UserID:      parsedID,
				Title:       reqBody.Title,
				Slug:        slug,
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
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			return c.Status(201).JSON(fiber.Map{
				"status":  "success",
				"message": "Project Added",
				"project": newProject,
			})
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	} else {
		return &fiber.Error{Code: 400, Message: "This title is not available."}
	}
}

func UpdateProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var project models.Project
	if err := initializers.DB.Where("id = ? AND user_id = ?", parsedProjectID, loggedInUserID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var reqBody schemas.ProjectUpdateSchema
	c.BodyParser(&reqBody)

	picName, err := utils.SaveFile(c, "coverPic", "project/coverPics", true, 1080, 1080)
	if err != nil {
		return err
	}
	reqBody.CoverPic = picName

	if reqBody.CoverPic != "" {
		err := utils.DeleteFile("project/coverPics", project.CoverPic)
		if err != nil {
			log.Printf("Error while deleting project cover pic: %e", err)
		}
	}

	projectValue := reflect.ValueOf(&project).Elem()
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Project updated successfully",
		"project": project,
	})
}

func DeleteProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id = ? AND user_id=?", parsedProjectID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var messages []models.Message
	if err := initializers.DB.Find(&messages, "project_id=?", parsedProjectID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}
	for _, message := range messages {
		if err := initializers.DB.Delete(&message).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}

	coverPic := project.CoverPic

	if err := initializers.DB.Delete(&project).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	err = utils.DeleteFile("project/coverPics", coverPic)
	if err != nil {
		log.Printf("Error while deleting project cover pic: %e", err)
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Project deleted successfully",
	})
}
