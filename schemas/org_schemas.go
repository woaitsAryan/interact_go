package schemas


/* Request body for creating a poll

used in controllers/organization_controllers/reviews.go
*/
type ReviewReqBody struct {
	ReviewContent string `json:"review_content"`
	ReviewRating  int    `json:"review_rating"`
	Anonymous     bool   `json:"anonymous"`
}
