package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddReport(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody schemas.ReportCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.ReportCreateSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	var existingReport models.Report
	if reqBody.UserID != "" {
		initializers.DB.Where("reporter_id=? AND user_id=?", parsedLoggedInUserID, reqBody.UserID).First(&existingReport)
	} else if reqBody.PostID != "" {
		initializers.DB.Where("reporter_id=? AND post_id=?", parsedLoggedInUserID, reqBody.PostID).First(&existingReport)
	} else if reqBody.ProjectID != "" {
		initializers.DB.Where("reporter_id=? AND project_id=?", parsedLoggedInUserID, reqBody.ProjectID).First(&existingReport)
	} else if reqBody.EventID != "" {
		initializers.DB.Where("reporter_id=? AND event_id=?", parsedLoggedInUserID, reqBody.EventID).First(&existingReport)
	} else if reqBody.OpeningID != "" {
		initializers.DB.Where("reporter_id=? AND opening_id=?", parsedLoggedInUserID, reqBody.OpeningID).First(&existingReport)
	} else if reqBody.GroupChatID != "" {
		initializers.DB.Where("reporter_id=? AND group_chat_id=?", parsedLoggedInUserID, reqBody.GroupChatID).First(&existingReport)
	} else if reqBody.ReviewID != "" {
		initializers.DB.Where("reporter_id=? AND review_id=?", parsedLoggedInUserID, reqBody.ReviewID).First(&existingReport)
	}

	if existingReport.ID != uuid.Nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "success",
			"message": "You have already filed a report.",
		})
	}

	report := models.Report{
		ReporterID: parsedLoggedInUserID,
		ReportType: reqBody.ReportType,
		Content:    reqBody.Content,
	}

	if reqBody.UserID != "" {
		parsedUserID, err := uuid.Parse(reqBody.UserID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid User ID"}
		}
		report.UserID = &parsedUserID
	} else if reqBody.PostID != "" {
		parsedPostID, err := uuid.Parse(reqBody.PostID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Post ID"}
		}
		report.PostID = &parsedPostID
	} else if reqBody.ProjectID != "" {
		parsedProjectID, err := uuid.Parse(reqBody.ProjectID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Project ID"}
		}
		report.ProjectID = &parsedProjectID
	} else if reqBody.EventID != "" {
		parsedEventID, err := uuid.Parse(reqBody.EventID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Event ID"}
		}
		report.EventID = &parsedEventID
	} else if reqBody.OpeningID != "" {
		parsedOpeningID, err := uuid.Parse(reqBody.OpeningID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Opening ID"}
		}
		report.OpeningID = &parsedOpeningID
	} else if reqBody.GroupChatID != "" {
		parsedGroupChatID, err := uuid.Parse(reqBody.GroupChatID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Group Chat ID"}
		}
		report.GroupChatID = &parsedGroupChatID
	} else if reqBody.ReviewID != "" {
		parsedReviewID, err := uuid.Parse(reqBody.ReviewID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Review ID"}
		}
		report.ReviewID = &parsedReviewID
	}

	result := initializers.DB.Create(&report)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	go routines.LogReport(&report)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Reported",
	})
}
