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

	report := models.Report{
		ReporterID: parsedLoggedInUserID,
		ReportType: reqBody.ReportType,
	}

	if reqBody.UserID != "" {
		parsedUserID, err := uuid.Parse(reqBody.UserID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid User ID"}
		}
		report.UserID = &parsedUserID
	}
	if reqBody.PostID != "" {
		parsedPostID, err := uuid.Parse(reqBody.PostID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Post ID"}
		}
		report.PostID = &parsedPostID
	}
	if reqBody.ProjectID != "" {
		parsedProjectID, err := uuid.Parse(reqBody.ProjectID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Project ID"}
		}
		report.ProjectID = &parsedProjectID
	}
	if reqBody.OpeningID != "" {
		parsedOpeningID, err := uuid.Parse(reqBody.OpeningID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Opening ID"}
		}
		report.OpeningID = &parsedOpeningID
	}

	result := initializers.DB.Create(&report)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	go routines.LogReport(&report)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Reported",
	})
}
