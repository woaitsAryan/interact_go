package organization_controllers

import (
	"errors"
	"time"

	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetOrganization(c *fiber.Ctx) error {
	orgID := c.Params("orgID") //TODO when the param id is not in uuid format and there is no check for parsedID, db throw error that invalid input format for id but response is Internal Server Error

	var organization models.Organization
	if err := initializers.DB.Preload("User").Preload("User.Profile").First(&organization, "id=?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Organization of this ID Found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"organization": organization,
	})
}

func GetOrganizationTasks(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var organization models.Organization
	if err := initializers.DB.
		Preload("Memberships").
		Preload("Memberships.User").
		Find(&organization, "id = ? ", orgID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var tasks []models.Task
	if err := initializers.DB.
		Preload("Users").
		Preload("SubTasks").
		Preload("SubTasks.Users").
		Find(&tasks, "organization_id = ? ", orgID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"tasks":        tasks,
		"organization": organization,
	})
}

func GetOrganizationChats(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	var organization models.Organization
	if err := initializers.DB.
		Preload("Memberships").
		Preload("Memberships.User").
		Find(&organization, "id = ? ", orgID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var chats []models.GroupChat
	if err := initializers.DB.
		Preload("User").
		Preload("Memberships").
		Preload("Memberships.User").
		Find(&chats, "organization_id = ? ", orgID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "",
		"chats":        chats,
		"organization": organization,
	})
}

func GetOrgEvents(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var events []models.Event //TODO add last edited n all fields
	if err := paginatedDB.
		Preload("Organization").
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"events":  events,
	})
}

func GetOrganizationHistory(c *fiber.Ctx) error {
	orgID := c.Params("orgID")

	paginatedDB := API.Paginator(c)(initializers.DB)
	var history []models.OrganizationHistory

	if err := paginatedDB.
		Preload("User").
		Preload("Post").
		Preload("Event").
		Preload("Project").
		Preload("Task").
		Preload("Invitation").
		Where("organization_id=?", orgID).
		Order("created_at DESC").
		Find(&history).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"history": history,
	})
}

func UpdateOrg(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	var organization models.Organization
	if err := initializers.DB.Preload("User").First(&organization, "id = ?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No organization of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	orgMemberID := c.GetRespHeader("orgMemberID")
	parsedOrgMemberID, err := uuid.Parse(orgMemberID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}

	user := organization.User
	var isOrgEdited = false

	var reqBody schemas.UserUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	// if err := helpers.Validate[schemas.UserUpdateSchema](reqBody); err != nil {
	// 	return &fiber.Error{Code: 400, Message: err.Error()}
	// }

	oldProfilePic := user.ProfilePic
	oldCoverPic := user.CoverPic

	picName, err := utils.UploadImage(c, "profilePic", helpers.UserProfileClient, 500, 500)
	if err != nil {
		return err
	}
	reqBody.ProfilePic = &picName

	coverName, err := utils.UploadImage(c, "coverPic", helpers.UserCoverClient, 900, 400)
	if err != nil {
		return err
	}
	reqBody.CoverPic = &coverName

	if reqBody.Name != nil && orgMemberID == c.GetRespHeader("loggedInUserID") {
		user.Name = *reqBody.Name
		organization.OrganizationTitle = *reqBody.Name
		isOrgEdited = true
	}
	if reqBody.Bio != nil {
		user.Bio = *reqBody.Bio
	}
	if reqBody.Tagline != nil {
		user.Tagline = *reqBody.Tagline
	}
	if reqBody.ProfilePic != nil && *reqBody.ProfilePic != "" {
		user.ProfilePic = *reqBody.ProfilePic
	}
	if reqBody.CoverPic != nil && *reqBody.CoverPic != "" {
		user.CoverPic = *reqBody.CoverPic
	}
	if reqBody.Tags != nil {
		user.Tags = *reqBody.Tags
	}
	if reqBody.Links != nil {
		user.Links = *reqBody.Links
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if isOrgEdited {
		if err := initializers.DB.Save(&organization).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}

	if *reqBody.ProfilePic != "" {
		go routines.DeleteFromBucket(helpers.UserProfileClient, oldProfilePic)
	}

	if *reqBody.CoverPic != "" {
		go routines.DeleteFromBucket(helpers.UserCoverClient, oldCoverPic)
	}

	parsedOrgID, err := uuid.Parse(orgID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}
	go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 14, nil, nil, nil, nil, nil, "")

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Organization updated successfully",
		"user":    user,
	})
}

func DeleteOrganization(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	userID := c.GetRespHeader("loggedInUserID")

	var reqBody struct {
		VerificationCode string `json:"otp"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id=?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No User of this ID Found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var organization models.Organization
	if err := initializers.DB.First(&organization, "id=?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No Organization of this ID Found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	data, err := cache.GetOtpFromCache(user.ID.String())
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data), []byte(reqBody.VerificationCode)); err != nil {
		return &fiber.Error{Code: 400, Message: "Incorrect OTP"}
	}

	organization.User.Active = false
	organization.User.DeactivatedAt = time.Now()

	if err := initializers.DB.Delete(&organization).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(202).JSON(fiber.Map{
		"status":  "success",
		"message": "Organization deleted successfully",
	})
}