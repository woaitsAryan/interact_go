package organization_controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAnnouncement(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	parsedAnnouncementID, err := uuid.Parse(c.Params("announcementID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Announcement ID."}
	}

	var announcement models.Announcement
	if err := initializers.DB.Preload("TaggedUsers").Where("id=?", parsedAnnouncementID).First(&announcement).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Announcement does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var organization models.Organization
	if err := initializers.DB.Preload("User").Preload("Memberships").First(&organization, "id = ?", announcement.OrganizationID).Error; err != nil {
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

	if !announcement.IsOpen && !isMember {
		return &fiber.Error{Code: 401, Message: "Cannot access this Announcement."}
	}

	announcement.Organization = organization

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"announcement": announcement,
	})
}

func GetOrgAnnouncements(c *fiber.Ctx) error {
	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var organization models.Organization
	if err := initializers.DB.Preload("User").Preload("Memberships").First(&organization, "id = ?", parsedOrgID).Error; err != nil {
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
	if err := paginatedDB.Preload("TaggedUsers").Where("organization_id = ?", parsedOrgID).Order("created_at DESC").Find(&announcements).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"organization":  organization,
		"announcements": announcements,
	})
}

func AddAnnouncement(c *fiber.Ctx) error {
	var reqBody schemas.AnnouncementCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedLoggedInUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))
	parsedOrgMemberID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	announcement := models.Announcement{
		OrganizationID: parsedOrgID,
		Title:          reqBody.Title,
		Content:        reqBody.Content,
		IsOpen:         reqBody.IsOpen,
	}

	var taggedUsers []models.User

	if reqBody.TaggedUsernames != nil {
		for _, username := range reqBody.TaggedUsernames {
			var user models.User
			if err := initializers.DB.First(&user, "username=?", username).Error; err == nil {
				taggedUsers = append(taggedUsers, user)
			}
		}

		announcement.TaggedUsers = taggedUsers
	}

	if err := initializers.DB.Create(&announcement).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if len(taggedUsers) > 0 {
		for _, user := range taggedUsers {
			go routines.SendTaggedNotification(user.ID, parsedLoggedInUserID, nil, &announcement.ID)
		}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 21, nil, nil, nil, nil, nil, nil, &announcement.ID, nil, nil, "")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":       "success",
		"message":      "Announcement added",
		"announcement": announcement,
	})
}

func EditAnnouncement(c *fiber.Ctx) error {
	var reqBody schemas.AnnouncementUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedLoggedInUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	parsedAnnouncementID, err := uuid.Parse(c.Params("announcementID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Announcement ID."}
	}

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	var announcement models.Announcement
	if err := initializers.DB.Where("id=? AND organization_id = ?", parsedAnnouncementID, parsedOrgID).First(&announcement).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Announcement does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.Title != "" {
		announcement.Title = reqBody.Title
	}
	if reqBody.Content != "" {
		announcement.Content = reqBody.Content
	}
	announcement.IsOpen = reqBody.IsOpen
	announcement.IsEdited = true

	if err := initializers.DB.Save(&announcement).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedLoggedInUserID, 23, nil, nil, nil, nil, nil, nil, &announcement.ID, nil, nil, "")

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "Announcement Edited",
		"announcement": announcement,
	})
}

func DeleteAnnouncement(c *fiber.Ctx) error {
	parsedAnnouncementID, err := uuid.Parse(c.Params("announcementID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Announcement ID."}
	}

	parsedLoggedInUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	var announcement models.Announcement
	if err := initializers.DB.Where("id=? AND organization_id = ?", parsedAnnouncementID, parsedOrgID).First(&announcement).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource Bucket does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete all the tagged users
	if err := tx.Model(&announcement).Association("TaggedUsers").Clear(); err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Delete(&announcement).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.MarkOrganizationHistory(parsedOrgID, parsedLoggedInUserID, 22, nil, nil, nil, nil, nil, nil, nil, nil, nil, announcement.Title)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Announcement deleted",
	})
}
