package select_fields

const (
	User         = "id, name, username, profile_pic, active"
	ShorterUser  = "id, name, username, active"
	ExtendedUser = "id, name, username, profile_pic, active, tagline, no_followers"
	Project      = "id, user_id, title, slug, cover_pic"
)
