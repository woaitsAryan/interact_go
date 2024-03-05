package project_controllers

import (
	"fmt"
	"reflect"

	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/controllers/auth_controllers"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetProject(c *fiber.Ctx) error {
	slug := c.Params("slug")
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	projectInCache, err := cache.GetProject(slug)
	if err == nil {
		go routines.UpdateProjectViews(projectInCache)
		if parsedLoggedInUserID != projectInCache.UserID {
			go routines.UpdateLastViewedProject(parsedLoggedInUserID, projectInCache.ID)
		}
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "",
			"project": projectInCache,
		})
	}

	var project models.Project
	if err := initializers.DB.Preload("User").Preload("Openings").Preload("Memberships").First(&project, "slug = ? AND is_private = ? ", slug, false).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var memberships []models.Membership
	if err := initializers.DB.Preload("User").Find(&memberships, "project_id = ?", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.UpdateProjectViews(&project)

	if parsedLoggedInUserID != project.UserID {
		go routines.UpdateLastViewedProject(parsedLoggedInUserID, project.ID)
	}

	_, count, err := utils.GetProjectViews(project.ID)
	if err != nil {
		return err
	}
	project.Views = count
	project.Memberships = memberships

	go cache.SetProject(project.Slug, &project)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"project": project,
	})
}

func GetWorkSpaceProject(c *fiber.Ctx) error {
	slug := c.Params("slug")

	projectInCache, err := cache.GetProject("-workspace--" + slug)
	if err == nil {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "",
			"project": projectInCache,
		})
	}

	var project models.Project
	if err := initializers.DB.Preload("User").Preload("Openings").First(&project, "slug = ?", slug).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var memberships []models.Membership
	if err := initializers.DB.Preload("User").Find(&memberships, "project_id = ?", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var invitations []models.Invitation
	if err := initializers.DB.Preload("User").Find(&invitations, "project_id = ? AND (status = 0 OR status = -1)", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	// var chats []models.GroupChat
	// if err := initializers.DB.Find(&chats, "project_id = ? ", project.ID).Error; err != nil {
	// 	return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage:Err.Error() ,Err: err}
	// }

	_, count, err := utils.GetProjectViews(project.ID)
	if err != nil {
		return err
	}
	project.Views = count
	project.Memberships = memberships
	project.Invitations = invitations //TODO remove if not required
	// project.Chats = chats //TODO only include chats you are part of

	go cache.SetProject("-workspace--"+project.Slug, &project)

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"project":      project,
		"privateLinks": project.PrivateLinks,
	})
}

func GetProjectHistory(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	paginatedDB := API.Paginator(c)(initializers.DB)
	var history []models.ProjectHistory

	if err := paginatedDB.
		Preload("Sender").
		Preload("User").
		Preload("Opening").
		Preload("Application").
		Preload("Invitation").
		Preload("Task").
		Where("project_id=?", projectID).
		Order("created_at DESC").
		Find(&history).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"history": history,
	})
}

func GetWorkSpaceProjectTasks(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var project models.Project
	if err := initializers.DB.First(&project, "slug = ? ", slug).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var tasks []models.Task
	if err := initializers.DB.
		Preload("Users").
		Find(&tasks, "project_id = ? ", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"tasks":   tasks,
	})
}

func GetWorkSpacePopulatedProjectTasks(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var project models.Project
	if err := initializers.DB.Preload("User").Preload("Memberships").Preload("Memberships.User").First(&project, "slug = ? ", slug).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var tasks []models.Task
	if err := initializers.DB.
		Preload("Users").
		Preload("SubTasks").
		Preload("SubTasks.Users").
		Find(&tasks, "project_id = ? ", project.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"tasks":   tasks,
		"project": project,
	})
}

func GetWorkSpaceProjectChats(c *fiber.Ctx) error {
	projectID := c.Params("projectID")

	var chats []models.GroupChat
	if err := initializers.DB.
		Preload("User").
		Preload("Memberships").
		Preload("Memberships.User").
		Find(&chats, "project_id = ? ", projectID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chats":   chats,
	})
}

func GetMyLikedProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var projectLikes []models.Like
	if err := initializers.DB.Where("user_id = ? AND project_id IS NOT NULL", loggedInUserID).Find(&projectLikes).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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

	paginatedDB := API.Paginator(c)(initializers.DB)

	var projects []models.Project
	if err := paginatedDB.
		Where("user_id = ? AND is_private = ?", userID, false).
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"projects": projects,
	})
}

func GetUserContributingProjects(c *fiber.Ctx) error {
	userID := c.Params("userID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var memberships []models.Membership
	if err := paginatedDB.
		Preload("Project").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&memberships).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var projects []models.Project
	for _, membership := range memberships {
		if !membership.Project.IsPrivate {
			projects = append(projects, membership.Project)
		}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	slug := slug.Make(reqBody.Title)
	suffix := 0
	newSlug := slug
	var existingProject models.Project
	for {
		if err := initializers.DB.Where("slug=?", newSlug).First(&existingProject).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				parsedID, err := uuid.Parse(c.GetRespHeader("loggedInUserID"))
				if err != nil {
					return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
				}

				var user models.User
				if err := initializers.DB.Where("id=?", parsedID).First(&user).Error; err != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
				}
				if !user.Verified {
					return &fiber.Error{Code: 401, Message: config.VERIFICATION_ERROR}
				}

				picName, err := utils.UploadImage(c, "coverPic", helpers.ProjectClient, 2560, 2560)
				if err != nil {
					return err
				}

				hash := "no-hash"

				if picName == "" {
					if defaultHash, found := config.AcceptedDefaultProjectHashes[reqBody.CoverPic]; found {
						picName = reqBody.CoverPic
						hash = defaultHash
					}
				}

				newProject := models.Project{
					UserID:      parsedID,
					Title:       reqBody.Title,
					Slug:        newSlug,
					Tagline:     reqBody.Tagline,
					CoverPic:    picName,
					BlurHash:    hash,
					Description: reqBody.Description,
					Tags:        reqBody.Tags,
					Category:    reqBody.Category,
					IsPrivate:   reqBody.IsPrivate,
					Links:       reqBody.Links,
				}

				orgMemberID := c.GetRespHeader("orgMemberID")
				orgID := c.Params("orgID")

				if orgMemberID != "" && orgID != "" {
					newProject.NumberOfMembers = 0
				}

				result := initializers.DB.Create(&newProject)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}

				if orgMemberID != "" && orgID != "" {
					parsedOrgMemberID, _ := uuid.Parse(orgMemberID)
					parsedOrgID, _ := uuid.Parse(orgID)
					go routines.MarkProjectHistory(newProject.ID, parsedOrgMemberID, -1, nil, nil, nil, nil, nil, "")
					go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 9, nil, &newProject.ID, nil, nil, nil, nil, nil, nil, nil, "")
					go routines.IncrementOrgProject(parsedOrgID)
				} else {
					go routines.MarkProjectHistory(newProject.ID, parsedID, -1, nil, nil, nil, nil, nil, "")
				}

				go routines.IncrementUserProject(parsedID)
				routines.GetImageBlurHash(c, "coverPic", &newProject)

				return c.Status(201).JSON(fiber.Map{
					"status":  "success",
					"message": "Project Added",
					"project": newProject,
				})
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		} else {
			suffix++
			newSlug = fmt.Sprintf("%s-%d", slug, suffix)
			existingProject = models.Project{}
		}
	}
}

func UpdateProject(c *fiber.Ctx) error {
	slug := c.Params("slug")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var project models.Project
	if err := initializers.DB.Where("slug = ? AND user_id = ?", slug, loggedInUserID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var reqBody schemas.ProjectUpdateSchema
	c.BodyParser(&reqBody)

	// picName, err := utils.SaveFile(c, "coverPic", "project/coverPics", true, 2560, 2560)
	picName, err := utils.UploadImage(c, "coverPic", helpers.ProjectClient, 2560, 2560)
	if err != nil {
		return err
	}
	reqBody.CoverPic = picName
	oldProjectPic := project.CoverPic

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

	if reqBody.IsPrivate {
		project.IsPrivate = true
	} else {
		project.IsPrivate = false
	}

	if err := initializers.DB.Save(&project).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.CoverPic != "" {
		go routines.DeleteFromBucket(helpers.ProjectClient, oldProjectPic)
	}

	projectMemberID := c.GetRespHeader("projectMemberID")
	parsedID, _ := uuid.Parse(projectMemberID)

	orgMemberID := c.GetRespHeader("orgMemberID")
	orgID := c.Params("orgID")
	if orgMemberID != "" {
		parsedOrgMemberID, _ := uuid.Parse(orgMemberID)
		parsedOrgID, _ := uuid.Parse(orgID)
		go routines.MarkProjectHistory(project.ID, parsedOrgMemberID, 2, nil, nil, nil, nil, nil, "")
		go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 11, nil, &project.ID, nil, nil, nil, nil, nil, nil, nil, "")
	} else {
		go routines.MarkProjectHistory(project.ID, parsedID, 2, nil, nil, nil, nil, nil, "")
	}

	routines.GetImageBlurHash(c, "coverPic", &project)

	go cache.RemoveProject(project.Slug)
	go cache.RemoveProject("-workspace--" + project.Slug)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Project updated successfully",
		"project": project,
	})
}

func DeleteProject(c *fiber.Ctx) error {
	actualLoggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedProjectUserID, _ := uuid.Parse(actualLoggedInUserID)
	if c.Params("orgID") != "" {
		actualLoggedInUserID = c.GetRespHeader("orgMemberID")
	}
	parsedLoggedInUserID, _ := uuid.Parse(actualLoggedInUserID)

	parsedProjectID, err := uuid.Parse(c.Params("projectID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var reqBody struct {
		VerificationCode string `json:"otp"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "OTP not provided."}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id = ? AND user_id=?", parsedProjectID, parsedProjectUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var user models.User
	if err := initializers.DB.Where("id=?", parsedLoggedInUserID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	data, err := cache.GetOtpFromCache(user.ID.String() + "-" + project.ID.String())
	if err != nil {
		return &fiber.Error{Code: 400, Message: "OTP Expired"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data), []byte(reqBody.VerificationCode)); err != nil {
		return &fiber.Error{Code: 400, Message: "Incorrect OTP"}
	}

	coverPic := project.CoverPic

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "DeleteProject")
		}
	}()

	var messages []models.Message
	if err := tx.Find(&messages, "project_id=?", parsedProjectID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}
	for _, message := range messages {
		if err := tx.Delete(&message).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	if err := tx.Delete(&project).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	orgMemberID := c.GetRespHeader("orgMemberID")
	orgID := c.Params("orgID")
	if orgMemberID != "" && orgID != "" {
		parsedOrgID, err := uuid.Parse(orgID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
		}

		parsedOrgMemberID, err := uuid.Parse(orgMemberID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid User ID."}
		}
		go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 10, nil, nil, nil, nil, nil, nil, nil, nil, nil, project.Title)
		go routines.DecrementOrgProject(parsedOrgID)
	}

	go routines.DeleteFromBucket(helpers.ProjectClient, coverPic)
	go routines.DecrementUserProject(parsedLoggedInUserID)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Project deleted successfully",
	})
}

func SendDeleteVerificationCode(c *fiber.Ctx) error {
	actualLoggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedProjectUserID, _ := uuid.Parse(actualLoggedInUserID)
	if c.Params("orgID") != "" {
		actualLoggedInUserID = c.GetRespHeader("orgMemberID")
	}
	parsedLoggedInUserID, _ := uuid.Parse(actualLoggedInUserID)

	parsedProjectID, err := uuid.Parse(c.Params("projectID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var user models.User
	if err := initializers.DB.Where("id=?", parsedLoggedInUserID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var project models.Project
	if err := initializers.DB.First(&project, "id = ? AND user_id=?", parsedProjectID, parsedProjectUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	code := auth_controllers.GenerateOTP(6)
	hash, err := bcrypt.GenerateFromPassword([]byte(code), 10)
	if err != nil {
		go helpers.LogServerError("Error while hashing an OTP.", err, c.Path())
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}
	err = helpers.SendMail(config.VERIFICATION_DELETE_PROJECT_SUBJECT, config.VERIFICATION_EMAIL_BODY+code, user.Name, user.Email, "<div><strong>This is Valid for next 10 minutes only!</strong></div>")
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	err = cache.SetOtpToCache(user.ID.String()+"-"+project.ID.String(), []byte(hash))
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "OTP sent to registered mail",
	})
}
