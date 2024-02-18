package utils

import "gorm.io/gorm"

type OrderType string

const (
	Latest     OrderType = "latest"
	MostLiked  OrderType = "most_liked"
	MostViewed OrderType = "most_viewed"
	Trending   OrderType = "trending"
)

type ModelType string

const (
	Projects ModelType = "projects"
	Users    ModelType = "users"
	Openings ModelType = "openings"
	Events   ModelType = "events"
	Posts    ModelType = "posts"
)

func Order(db *gorm.DB, order OrderType, modelType ModelType) *gorm.DB {
	switch order {
	case Latest:
		return db.Order(string(modelType) + ".created_at DESC")

	case MostLiked:
		if modelType != Users {
			return db.Order(string(modelType) + ".no_likes DESC")
		}
		return db

	case MostViewed:
		return db.Order(string(modelType) + ".impressions DESC")

	case "", Trending:
		switch modelType {
		case Projects:
			return db.Where("is_private = ?", false).
				Select("*, (total_no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
				Order("weighted_average DESC")

		case Users:
			return db.Where("active=?", true).
				Where("organization_status=? AND verified=? AND username != users.email", false, true).
				Omit("users.phone_no").
				Omit("users.email").
				Select("*, (0.6 * no_followers - 0.4 * no_following + 0.3 * total_no_views) / (1 + EXTRACT(EPOCH FROM age(NOW(), created_at)) / 3600 / 24 / 21) AS weighted_average"). //* 21 days
				Order("weighted_average DESC, created_at ASC")

		case Openings:
			return db.Where("openings.active=true").
				Joins("JOIN openings ON openings.project_id = projects.id AND projects.is_private = ?", false).
				Select("openings.*, (projects.total_no_views * 0.5 + openings.no_of_applications * 0.3) / (1 + EXTRACT(EPOCH FROM age(NOW(), openings.created_at)) / 3600 / 24 / 15) AS t_ratio").
				Order("t_ratio DESC")

		case Events:
			return db.Order(string(modelType) + ".impressions DESC")

		case Posts:
			return db.Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
				Select("*, posts.id, posts.created_at, (2 * no_likes + no_comments + 5 * no_shares) / (1 + EXTRACT(EPOCH FROM age(NOW(), posts.created_at)) / 3600 / 24 / 7) AS weighted_average"). //* 7 days
				Order("weighted_average DESC, posts.created_at ASC")
		default:
			return db.Order(string(modelType) + ".created_at DESC")
		}

	default:
		return db
	}
}
