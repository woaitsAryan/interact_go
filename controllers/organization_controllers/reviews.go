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
	var reqBody schemas.ReviewReqBody
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

	var review models.OrganizationReview
	if err := initializers.DB.Where("user_id = ? AND organization_id = ?", parsedUserID, parsedOrgID).First(&review).Error; err == nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Review already exists."}
	} else {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error creating review."}
		}
	}

	review = models.OrganizationReview{
		UserID:         parsedUserID,
		OrganizationID: parsedOrgID,
		Review:         reqBody.ReviewContent,
		Rating:         reqBody.ReviewRating,
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

/* Deletes an organizational review

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

	var review models.OrganizationReview
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

/* Fetches all organizational reviews

Takes input of organization ID to fetch all reviews for an organization.
If anonymous is true then omits user information
*/
func FetchReviews(c *fiber.Ctx) error {

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	var reviews []models.OrganizationReview
	if err := initializers.DB.Where("organization_id = ?", parsedOrgID).
		Order("like_count desc").
		Omit("relevance").
		Find(&reviews).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusNotFound, Message: "No reviews found."}
		}
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error fetching reviews."}
	}
	for i := range reviews {
		if reviews[i].Anonymous {
			reviews[i].UserID = uuid.Nil
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"reviews": reviews,
	})
}

/* Likes an organizational review

Takes input of review ID to like a review.
Currently an user can like their own review.
Also checks so you can't like a review multiple times
*/
func LikeReview(c *fiber.Ctx) error {
	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid User ID."}
	}
	reviewID, err := uuid.Parse(c.Params("reviewID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Review ID."}
	}
	var review models.OrganizationReview
	if err := initializers.DB.Where("id = ?", reviewID).First(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	for _, id := range review.LikedBy {
		if id == parsedUserID.String() {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User has already liked this review."}
		}
	}

	review.LikeCount++
	review.LikedBy = append(review.LikedBy, parsedUserID.String())
	if err := initializers.DB.Save(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error liking review."}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Review liked",
	})
}

/* Remove like from an organizational review

Takes input of review ID to remove like from a review.
Checks if the user has liked to be able to unlike.
*/
func RemoveLike(c *fiber.Ctx) error {
	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid User ID."}
	}
	reviewID, err := uuid.Parse(c.Params("reviewID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Review ID."}
	}
	var review models.OrganizationReview
	if err := initializers.DB.Where("id = ?", reviewID).First(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	isFound := false

	for i, id := range review.LikedBy {
		if id == parsedUserID.String() {
			review.LikedBy = append(review.LikedBy[:i], review.LikedBy[i+1:]...)
			review.LikeCount--
			isFound = true
			break
		}
	}
	if !isFound {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User has not liked this review."}
	}
	if err := initializers.DB.Save(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error removing like from review."}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Like removed from review",
	})
}

/* Dislike a review

Takes input of review ID to dislike a review.
Currently an user can dislike their own review.
Also checks so you can't dislike a review multiple times
*/
func DislikeReview(c *fiber.Ctx) error {
	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid User ID."}
	}
	reviewID, err := uuid.Parse(c.Params("reviewID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Review ID."}
	}
	var review models.OrganizationReview
	if err := initializers.DB.Where("id = ?", reviewID).First(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	for _, id := range review.DislikedBy {
		if id == parsedUserID.String() {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User has already disliked this review."}
		}
	}

	review.DislikeCount++
	review.DislikedBy = append(review.DislikedBy, parsedUserID.String())
	if err := initializers.DB.Save(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error disliking review."}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Review disliked",
	})
}

/* Remove dislike from a review

Takes input of review ID to remove dislike from a review.
Checks if the user has disliked to be able to undislike.
*/
func RemoveDislike(c *fiber.Ctx) error {
	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid User ID."}
	}
	reviewID, err := uuid.Parse(c.Params("reviewID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Review ID."}
	}
	var review models.OrganizationReview
	if err := initializers.DB.Where("id = ?", reviewID).First(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error fetching reviews"}
	}

	isFound := false

	for i, id := range review.DislikedBy {
		if id == parsedUserID.String() {
			review.DislikedBy = append(review.DislikedBy[:i], review.DislikedBy[i+1:]...)
			review.DislikeCount--
			isFound = true
			break
		}
	}
	if !isFound {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User has not disliked this review."}
	}
	if err := initializers.DB.Save(&review).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Error removing dislike from review."}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Dislike removed from review",
	})
}
