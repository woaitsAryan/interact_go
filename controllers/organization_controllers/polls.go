package organization_controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
		Question:       reqBody.Question,
		IsMultiAnswer:  reqBody.IsMultiAnswer,
	}

	if err := initializers.DB.Create(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	tx := initializers.DB.Begin()

	for _, optionText := range reqBody.Options {
		option := &models.Option{
			PollID: poll.ID,
			Text:   optionText,
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
	
	go routines.MarkOrganizationHistory(orgID, parsedUserID, 18, nil, nil, nil, nil, nil, &poll.ID, "")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll created!",
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
	if err := initializers.DB.First(&poll, "id = ?", parsedPollID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if !poll.IsMultiAnswer {
		for _, option := range poll.Options {
			for _, voter := range option.VotedBy {
				if voter.ID == parsedUserID {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User has already voted"})
				}
			}
		}
	}

	var option models.Option
	if err := initializers.DB.First(&option, "id = ?", parsedOptionID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", parsedUserID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	option.Votes++
	option.VotedBy = append(option.VotedBy, user)

	if err := initializers.DB.Save(&option).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
	for i, voter := range option.VotedBy {
		if voter.ID == parsedUserID {
			option.VotedBy = append(option.VotedBy[:i], option.VotedBy[i+1:]...)
			option.Votes--
			break
		}
	}

	if err := initializers.DB.Model(&option).Association("VotedBy").Replace(option.VotedBy); err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Save(&option).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Vote removed!",
	})
}

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

	var polls []models.Poll
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	if err := initializers.DB.Preload("Options").Preload("Options.VotedBy").Where("organization_id = ? AND created_at >= ?", orgID, oneWeekAgo).Find(&polls).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"poll":   polls,
	})
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
	if err := initializers.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	pollOption := models.Poll{ID: pollID}
	if err := initializers.DB.Model(&pollOption).Association("Options").Delete(); err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	orgID := poll.OrganizationID
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("orgMemberID"))

	go routines.MarkOrganizationHistory(orgID, parsedUserID, 19, nil, nil, nil, nil, nil, nil, poll.Question)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll deleted!",
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

	poll.Question = reqBody.Question
	poll.IsEdited = true

	if err := initializers.DB.Save(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	orgID := poll.OrganizationID

	go routines.MarkOrganizationHistory(orgID, parsedUserID, 20, nil, nil, nil, nil, nil, &poll.ID, "")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll edited!",
	})
}