package schemas

type ReportCreateSchema struct {
	Content    string `json:"content" validate:"required,max=500"`
	ReportType int    `json:"reportType" validate:"required"`
	UserID     string `json:"userID"`
	PostID     string `json:"postID"`
	ProjectID  string `json:"projectID"`
	OpeningID  string `json:"openingID"`
}
