package organization_controllers

import (
	"errors"
	"sort"
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var tasks []models.Task
	if err := initializers.DB.
		Preload("Users").
		Preload("SubTasks").
		Preload("SubTasks.Users").
		Find(&tasks, "organization_id = ? ", orgID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var chats []models.GroupChat
	if err := initializers.DB.
		Preload("User").
		Preload("Memberships").
		Preload("Memberships.User").
		Order("created_at DESC").
		Find(&chats, "organization_id = ? ", orgID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		Preload("Coordinators").
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		Preload("Announcement").
		Preload("Invitation").
		Preload("Invitation.User").
		Where("organization_id=?", orgID).
		Order("created_at DESC").
		Find(&history).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"history": history,
	})
}

type NewsFeedItem interface {
	GetCreatedAt() time.Time
}

// Announcement implements NewsFeedItem interface
type AnnouncementAlias models.Announcement

// GetCreatedAt is a method for AnnouncementAlias to satisfy the interface
func (a AnnouncementAlias) GetCreatedAt() time.Time {
	return a.CreatedAt
}

// Poll implements NewsFeedItem interface
type PollAlias models.Poll

// GetCreatedAt is a method for PollAlias to satisfy the interface
func (p PollAlias) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func GetOrgNewsFeed(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid organization ID."}
	}

	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var organization models.Organization
	if err := initializers.DB.Preload("User").Preload("Memberships").First(&organization, "id = ?", orgID).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid organization ID."}
	}

	isMember := false
	if organization.UserID == parsedUserID {
		isMember = true
	} else {
		for _, membership := range organization.Memberships {
			if membership.UserID == parsedUserID {
				isMember = true
				break
			}
		}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	if !isMember {
		paginatedDB = paginatedDB.Where("is_open = ?", true)
	}

	var announcements []models.Announcement
	if err := paginatedDB.Preload("TaggedUsers").Where("organization_id = ?", orgID).Order("created_at DESC").Find(&announcements).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	paginatedDB = API.Paginator(c)(initializers.DB)

	if !isMember {
		paginatedDB = paginatedDB.Where("is_open = ?", true)
	}

	db := paginatedDB.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("options.created_at DESC")
	}).Preload("Options.VotedBy", LimitedUsers).Where("organization_id = ?", orgID)

	var polls []models.Poll
	if err := db.Order("created_at DESC").Find(&polls).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	// Combine announcements and polls, and sort them by created_at
	var newsFeed []NewsFeedItem
	for _, a := range announcements {
		newsFeed = append(newsFeed, AnnouncementAlias(a))
	}
	for _, p := range polls {
		newsFeed = append(newsFeed, PollAlias(p))
	}

	// Sort the combined news feed by created_at
	sort.Slice(newsFeed, func(i, j int) bool {
		return newsFeed[i].GetCreatedAt().After(newsFeed[j].GetCreatedAt())
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"newsFeed":     newsFeed,
		"organization": organization,
	})
}

func UpdateOrg(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	var organization models.Organization
	if err := initializers.DB.Preload("User").First(&organization, "id = ?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No organization of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if isOrgEdited {
		if err := initializers.DB.Save(&organization).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
	go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 14, nil, nil, nil, nil, nil, nil, nil, nil, "")

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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var organization models.Organization
	if err := initializers.DB.First(&organization, "id=?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No Organization of this ID Found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Organization deleted successfully",
	})
}
