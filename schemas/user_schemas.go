package schemas

type UserCreateSchema struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"alphanum,required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type UserUpdateSchema struct {
	Name       string   `json:"name" validate:"alpha"`
	ProfilePic string   `json:"profilePic" validate:"image"`
	CoverPic   string   `json:"coverPic" validate:"image"`
	Bio        string   `json:"bio"`
	Title      string   `json:"title"`
	Tagline    string   `json:"tagline"`
	Tags       []string `json:"tags" validate:"dive,alpha"`
}

type AchievementCreateSchema struct {
	Achievements []AchievementSchema `json:"achievements"`
}

type AchievementSchema struct {
	ID     string   `json:"id"`
	Title  string   `json:"title" validate:"alpha"`
	Skills []string `json:"skills" validate:"dive,alpha"`
}
