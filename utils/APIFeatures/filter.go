package utils

import (
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Filter(c *fiber.Ctx, index int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		order := c.Query("order", "")
		switch index {
		//* For Project
		case 1:
			projectFields := []string{"tags", "category"}
			for _, projectField := range projectFields {
				value := c.Query(projectField, "")
				db = genericFilter(db, projectField, value, "projects")
			}

			return getOrderedDB(db, order, "projects")
		//* For User
		case 2:
			userFields := []string{"tags", "location", "areas_of_collaboration", "school"}
			for _, userField := range userFields {
				value := c.Query(userField, "")
				db = genericFilter(db, userField, value, "users")
			}

			return getOrderedDB(db, order, "users")
		//* For Event
		case 3:
			eventFields := []string{"tags", "location"}
			for _, eventField := range eventFields {
				value := c.Query(eventField, "")
				db = genericFilter(db, eventField, value, "events")
			}

			db = eventTimeSearch(c, db)
			return getOrderedDB(db, order, "events")
		//* For Opening
		case 4:
			openingFields := []string{"tags"}
			for _, openingField := range openingFields {
				value := c.Query(openingField, "")
				db = genericFilter(db, openingField, value, "openings")
			}

			return getOrderedDB(db, order, "openings")
		default:
			return db
		}
	}
}

func getOrderedDB(db *gorm.DB, order, modelType string) *gorm.DB {
	switch order {
	case "latest":
		return db.Order(modelType + ".created_at DESC")

	case "most_liked":
		if modelType != "users" {
			return db.Order(modelType + ".no_likes DESC")
		}
		return db

	case "most_viewed":
		return db.Order(modelType + ".impressions DESC")

	case "", "trending":
		switch modelType {
		case "projects":
			return db.Where("is_private = ?", false).
				Select("*, (total_no_views + 3 * no_likes + 2 * no_comments + 5 * no_shares) AS weighted_average").
				Order("weighted_average DESC")

		case "users":
			return db.Where("active=?", true).
				Where("organization_status=? AND verified=? AND username != users.email", false, true).
				Omit("users.phone_no").
				Omit("users.email").
				Select("*, (0.6 * no_followers - 0.4 * no_following + 0.3 * total_no_views) / (1 + EXTRACT(EPOCH FROM age(NOW(), created_at)) / 3600 / 24 / 21) AS weighted_average"). //! 21 days
				Order("weighted_average DESC, created_at ASC")

		case "openings":
			return db.Where("openings.active=true").
				Joins("JOIN openings ON openings.project_id = projects.id AND projects.is_private = ?", false).
				Select("openings.*, (projects.total_no_views * 0.5 + openings.no_of_applications * 0.3) / (1 + EXTRACT(EPOCH FROM age(NOW(), openings.created_at)) / 3600 / 24 / 15) AS t_ratio").
				Order("t_ratio DESC")

		case "events":
			return db.Order(modelType + ".impressions DESC")

		case "posts":
			return db.Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
				Select("*, posts.id, posts.created_at, (2 * no_likes + no_comments + 5 * no_shares) / (1 + EXTRACT(EPOCH FROM age(NOW(), posts.created_at)) / 3600 / 24 / 7) AS weighted_average"). //! 7 days
				Order("weighted_average DESC, posts.created_at ASC")
		default:
			return db.Order(modelType + ".created_at DESC")
		}

	default:
		return db
	}
}

func eventTimeSearch(c *fiber.Ctx, db *gorm.DB) *gorm.DB {
	//* Get Events Between start and end

	start := c.Query("start", "")
	end := c.Query("end", "")

	if start != "" && end != "" {
		startTime, err := time.Parse(time.RFC3339, start)
		if err != nil {
			helpers.LogServerError("Error parsing start timestamp", err, "timestampSearch")
			return db
		}

		endTime, err := time.Parse(time.RFC3339, end)
		if err != nil {
			helpers.LogServerError("Error parsing end timestamp", err, "timestampSearch")
			return db
		}

		return db.Where("start_time BETWEEN ? AND ?", startTime, endTime)
	} else if start != "" {
		startTime, err := time.Parse(time.RFC3339, start)
		if err != nil {
			helpers.LogServerError("Error parsing start timestamp", err, "timestampSearch")
			return db
		}

		return db.Where("start_time >= ?", startTime)
	} else if end != "" {
		endTime, err := time.Parse(time.RFC3339, end)
		if err != nil {
			helpers.LogServerError("Error parsing end timestamp", err, "timestampSearch")
			return db
		}

		return db.Where("start_time <= ?", endTime)
	}
	return db
}

func genericFilter(db *gorm.DB, field, value string, modelType string) *gorm.DB {
	if value != "" {
		if isArrayField(field) {
			db = db.Where("? =  ANY("+modelType+"."+field+")", value)
		} else if isProfile(field) {
			if field == "areas_of_collaboration" {
				db = db.Joins("JOIN profiles ON users.id = profiles.user_id").Where("? =  ANY("+"profiles."+field+")", value)
			} else {
				db = db.Joins("JOIN profiles ON users.id = profiles.user_id").Where("profiles."+field+" ILIKE ?", "%"+value+"%")
			}
		} else {
			db = db.Where(modelType+"."+field+" ILIKE ?", "%"+value+"%")
		}
	}
	return db
}

func isArrayField(field string) bool {
	switch field {
	case "tags":
		return true
	default:
		return false
	}
}

func isProfile(field string) bool {
	switch field {
	case "location", "areas_of_collaboration", "school":
		return true
	default:
		return false
	}
}
