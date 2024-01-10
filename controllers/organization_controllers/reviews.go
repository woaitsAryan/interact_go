package organization_controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
	Adds an organizational review

Takes input of content, rating and anonymity to make a review for an organization.
Has a go routine to compute relevance of the review.
*/
func AddReview(c *fiber.Ctx) error {
	//TODO switch to AppError where required
	var reqBody schemas.ReviewCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid User ID."}
	}
	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	var review models.Review
	if err := initializers.DB.Where("user_id = ? AND organization_id = ?", parsedUserID, parsedOrgID).First(&review).Error; err == nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Review already exists."}
	} else {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error creating review."}
		}
	}

	review = models.Review{
		UserID:         parsedUserID,
		OrganizationID: parsedOrgID,
		Content:        reqBody.Content,
		Rating:         reqBody.Rating,
		Anonymous:      reqBody.Anonymous,
	}
	if err := initializers.DB.Create(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error creating review."}
	}

	go routines.ComputeRelevance(review.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Review added",
	})
}

/*
	Deletes an organizational review

Takes input of organization ID and user ID to delete a review for an organization.
*/
func DeleteReview(c *fiber.Ctx) error {
	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid User ID."}
	}
	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	var review models.Review
	if err := initializers.DB.Where("user_id = ? AND organization_id = ?", parsedUserID, parsedOrgID).First(&review).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Review does not exist."}
		}
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error deleting review."}
	}
	if err := initializers.DB.Delete(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error deleting review."}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Review deleted",
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

	var reviews []models.Review
	if err := initializers.DB.Where("organization_id = ?", parsedOrgID).
		Order("number_of_up_votes desc").
		Find(&reviews).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error fetching reviews."}
	}

	// for i := range reviews {
	// 	if reviews[i].Anonymous {
	// 		reviews[i].UserID = uuid.Nil
	// 		reviews[i].User = models.User{}
	// 	}
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"reviews": reviews,
	})
}
