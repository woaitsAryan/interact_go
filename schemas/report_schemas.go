package schemas

type ReportCreateSchema struct {
	Content     string `json:"content" validate:"max=1000"`
	ReportType  int    `json:"reportType" validate:"required"`
	UserID      string `json:"userID"`
	PostID      string `json:"postID"`
	ProjectID   string `json:"projectID"`
	EventID     string `json:"eventID"`
	OpeningID   string `json:"openingID"`
	GroupChatID string `json:"chatID"`
}
