package organization_controllers

import (
	"math"

	"github.com/Pratham-Mishra04/interact/cache"
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

type ReviewData struct {
	Total   int          `json:"total"`
	Average float64      `json:"average"`
	Counts  map[int8]int `json:"counts"`
}

func GetOrgReviewData(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	parsedOrgID, err := uuid.Parse(orgID)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	var reviewData ReviewData
	if err = cache.GetFromCacheGeneric("review-data-"+orgID, reviewData); err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "success",
			"reviewData": reviewData,
		})
	}

	var reviews []models.Review
	if err := initializers.DB.Where("organization_id = ?", parsedOrgID).Find(&reviews).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	ratingCounts := make(map[int8]int)

	totalRatings := 0
	avgRating := 0.0

	ratingCounts[1] = 0
	ratingCounts[2] = 0
	ratingCounts[3] = 0
	ratingCounts[4] = 0
	ratingCounts[5] = 0

	for _, review := range reviews {
		ratingCounts[review.Rating]++
		totalRatings += int(review.Rating)
	}

	if len(reviews) > 0 {
		avgRating = float64(totalRatings) / float64(len(reviews))
		avgRating = math.Round(avgRating*100) / 100 // Round to two decimal places
	}

	reviewData.Average = avgRating
	reviewData.Counts = ratingCounts
	reviewData.Total = len(reviews)

	go cache.SetToCacheGeneric("review-data-"+orgID, reviewData)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"reviewData": reviewData,
	})
}

/*
	Fetches all organizational reviews

Takes input of organization ID to fetch all reviews for an organization.
If anonymous is true then omits user information
*/
func FetchOrgReviews(c *fiber.Ctx) error {
	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var reviews []models.Review
	if err := paginatedDB.Where("organization_id = ?", parsedOrgID).
		Preload("User").
		Order("number_of_up_votes DESC, created_at DESC").
		Find(&reviews).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"reviews": reviews,
	})
}

/*
	Adds an organizational review

Takes input of content, rating and anonymity to make a review for an organization.
Has a go routine to compute relevance of the review.
*/
func AddReview(c *fiber.Ctx) error {
	var reqBody schemas.ReviewCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid User ID."}
	}

	orgID := c.Params("orgID")
	parsedOrgID, err := uuid.Parse(orgID)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	review := models.Review{
		UserID:         &parsedUserID,
		OrganizationID: parsedOrgID,
		Content:        reqBody.Content,
		Rating:         reqBody.Rating,
		Anonymous:      reqBody.Anonymous,
	}
	if reqBody.Anonymous {
		review.UserID = nil
	}

	if err := initializers.DB.Create(&review).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.ComputeRelevance(review.ID)
	go cache.RemoveFromCacheGeneric("review-data-" + orgID)

	if err := initializers.DB.Preload("User").First(&review).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Review added",
		"review":  review,
	})
}

/*
	Deletes an organizational review

Takes input of organization ID and user ID to delete a review for an organization.
*/
func DeleteReview(c *fiber.Ctx) error {
	parsedReviewID, err := uuid.Parse(c.Params("reviewID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Review ID."}
	}

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid User ID."}
	}

	orgID := c.Params("orgID")
	parsedOrgID, err := uuid.Parse(orgID)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	var review models.Review
	if err := initializers.DB.Where("id=? AND user_id = ? AND organization_id = ?", parsedReviewID, parsedUserID, parsedOrgID).First(&review).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Review does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&review).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go cache.RemoveFromCacheGeneric("review-data-" + orgID)

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Review deleted",
	})
}
