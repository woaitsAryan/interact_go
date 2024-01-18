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

/*
	Fetches all polls in an organization.

Reads the organization ID from request params
Fetches all polls created in the last week
*/
func FetchPolls(c *fiber.Ctx) error {
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
	// Check if the user is a member of the organization
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

	// If the user is not a member, only show open polls
	if !isMember {
		paginatedDB = paginatedDB.Where("is_open = ?", true)
	}

	// oneWeekAgo := time.Now().AddDate(0, 0, -7)
	// db := initializers.DB.Preload("Options").Preload("Options.VotedBy", LimitedUsers).Where("organization_id = ? AND created_at >= ?", orgID, oneWeekAgo)
	db := paginatedDB.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("options.created_at DESC")
	}).Preload("Options.VotedBy", LimitedUsers).Where("organization_id = ?", orgID)

	var polls []models.Poll
	if err := db.Order("created_at DESC").Find(&polls).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"polls":        polls,
		"organization": organization,
	})
}

/*
	Creates a new poll.

It reads the poll data from the request body
If the request body is invalid, it returns a 400 status code.
If the user ID is invalid, it returns a 500 status code.
*/
func CreatePoll(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))
	var reqBody schemas.CreatePollRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}
	if len(reqBody.Options) < 2 || len(reqBody.Options) > 10 {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}

	orgID, _ := uuid.Parse(c.Params("orgID"))

	var poll = models.Poll{
		OrganizationID: orgID,
		Title:          reqBody.Title,
		Content:        reqBody.Content,
		IsMultiAnswer:  reqBody.IsMultiAnswer,
		IsOpen:         reqBody.IsOpen,
	}

	if err := initializers.DB.Create(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	tx := initializers.DB.Begin()

	for _, optionText := range reqBody.Options {
		option := &models.Option{
			PollID:  poll.ID,
			Content: optionText,
		}
		if err := tx.Create(&option).Error; err != nil {
			tx.Rollback()
			return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.MarkOrganizationHistory(orgID, parsedUserID, 18, nil, nil, nil, nil, nil, &poll.ID, nil, "")

	if err := initializers.DB.Preload("Options").First(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll created!",
		"poll":    poll,
	})
}

/*
	Vote for an option in a poll.

It reads the OptionID and PollID from the request params
For a single answer poll, it checks if the user has already voted
Then it increments the vote count for the option and adds the user to the votedBy array
*/
func VotePoll(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	parsedPollID, err := uuid.Parse(c.Params("pollID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Poll ID."}
	}

	parsedOptionID, err := uuid.Parse(c.Params("OptionID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Option ID."}
	}

	var poll models.Poll
	if err := initializers.DB.Preload("Options").Preload("Options.VotedBy").First(&poll, "id = ?", parsedPollID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	votedOptionID := uuid.Nil

	for _, option := range poll.Options {
		for _, voter := range option.VotedBy {
			if voter.ID == parsedUserID {
				votedOptionID = option.ID
				if !poll.IsMultiAnswer {
					return &fiber.Error{Code: fiber.StatusBadRequest, Message: "You have already voted"}
				}
			}
		}
	}

	var option models.Option
	if err := initializers.DB.First(&option, "id = ?", parsedOptionID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if votedOptionID == option.ID {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Vote recorded!",
		})
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", parsedUserID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	option.Votes++
	poll.TotalVotes++
	option.VotedBy = append(option.VotedBy, user)

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "VotePoll")
		}
	}()

	if err := tx.Save(&option).Error; err != nil {
		return err
	}

	if err := tx.Save(&poll).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Vote recorded!",
	})
}

/*
	Remove a vote for an option in a poll.

Reads the OptionID from request params
Removes the user from the votedBy array and decrements the vote count
*/
func UnvotePoll(c *fiber.Ctx) error {

	parsedUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	parsedOptionID, err := uuid.Parse(c.Params("OptionID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Option ID."}
	}

	var option models.Option
	if err := initializers.DB.Preload("VotedBy").First(&option, "id = ?", parsedOptionID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var poll models.Poll
	if err := initializers.DB.First(&poll, "id = ?", option.PollID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	isFound := false

	for i, voter := range option.VotedBy {
		if voter.ID == parsedUserID {
			option.VotedBy = append(option.VotedBy[:i], option.VotedBy[i+1:]...)
			option.Votes--
			poll.TotalVotes--
			isFound = true
			break
		}
	}

	if !isFound {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User has not voted"}
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "UnVotePoll")
		}
	}()

	if err := tx.Model(&option).Association("VotedBy").Replace(option.VotedBy); err != nil {
		return err
	}

	if err := tx.Save(&option).Error; err != nil {
		return err
	}

	if err := tx.Save(&poll).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Vote removed!",
	})
}

/*	Edit a poll.

Reads the poll ID from request params
Reads the new question from request body
*/

func EditPoll(c *fiber.Ctx) error {
	var reqBody schemas.EditPollRequest
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}

	pollID, _ := uuid.Parse(c.Params("pollID"))

	var poll models.Poll
	if err := initializers.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	poll.Content = reqBody.Content
	poll.IsEdited = true
	poll.IsOpen = reqBody.IsOpen

	if err := initializers.DB.Save(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	orgID := poll.OrganizationID

	go routines.MarkOrganizationHistory(orgID, parsedUserID, 20, nil, nil, nil, nil, nil, &poll.ID, nil, "")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll edited!",
	})
}

func LimitedUsers(db *gorm.DB) *gorm.DB {
	return db.Limit(3)
}

/*
	Deletes a poll.

Reads the poll ID from request params
Deletes the poll cascading all the options
*/
func DeletePoll(c *fiber.Ctx) error {
	pollID, err := uuid.Parse(c.Params("pollID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid poll ID."}
	}

	var poll models.Poll
	if err := initializers.DB.Preload("Options").First(&poll, "id = ?", pollID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No Poll if this ID Found."}
		}
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	tx := initializers.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, option := range poll.Options {
		if err := tx.Model(&option).Association("VotedBy").Clear(); err != nil {
			tx.Rollback()
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	orgID := poll.OrganizationID

	if err := tx.Delete(&poll).Error; err != nil {
		tx.Rollback()
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	parsedUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	go routines.MarkOrganizationHistory(orgID, parsedUserID, 19, nil, nil, nil, nil, nil, nil, nil, poll.Content)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll deleted!",
	})
}
