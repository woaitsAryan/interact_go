package routines

import (
	"math/rand"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

/* Routine to compute relevance of a review 

used in controllers/organization_controllers/reviews.go
*/
func ComputeRelevance(reviewID uuid.UUID) {
	var review models.OrganizationReview
	if err := initializers.DB.First(&review, "id=?", reviewID).Error; err != nil {
		helpers.LogDatabaseError("No Review of this ID found-ComputeRelevance.", err, "go_routine")
	}
	review.Relevance = rand.Intn(91) + 10
	if err := initializers.DB.Save(&review).Error; err != nil {
		helpers.LogDatabaseError("Error in saving Review-ComputeRelevance.", err, "go_routine")
	}
}