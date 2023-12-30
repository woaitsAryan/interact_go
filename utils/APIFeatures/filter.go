package utils

import (
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Filter(c *fiber.Ctx, index int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// order := c.Query("order", "")
		switch index {
		//* For Project
		case 1:
			projectFields := []string{"tags", "category"}
			for _, projectField := range projectFields {
				value := c.Query(projectField, "")
				db = genericFilter(db, projectField, value, "projects")
			}

			return db
		//* For User
		case 2:
			userFields := []string{"tags", "location", "areas_of_collaboration", "school"}
			for _, userField := range userFields {
				value := c.Query(userField, "")
				db = genericFilter(db, userField, value, "users")
			}

			return db
		//* For Event
		case 3:
			eventFields := []string{"tags", "location"}
			for _, eventField := range eventFields {
				value := c.Query(eventField, "")
				db = genericFilter(db, eventField, value, "events")
			}

			return eventTimeSearch(c, db)
		//* For Opening
		case 4:
			openingFields := []string{"tags"}
			for _, openingField := range openingFields {
				value := c.Query(openingField, "")
				db = genericFilter(db, openingField, value, "openings")
			}

			return db

		//* For Tasks
		case 5:
			taskFields := []string{"tags", "priority", "is_completed", "deadline"}
			//TODO assigned to a user filter
			for _, taskField := range taskFields {
				value := c.Query(taskField, "")
				db = genericFilter(db, taskField, value, "tasks")
			}

			return db

			//* For Sub Tasks
		case 6:
			taskFields := []string{"tags", "priority", "is_completed", "deadline"}
			//TODO assigned to a user filter
			for _, taskField := range taskFields {
				value := c.Query(taskField, "")
				db = genericFilter(db, taskField, value, "sub_tasks")
			}

			return db
		default:
			return db
		}
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
		} else if isBooleanField(field) {
			if field == "true" {
				db = db.Where(modelType+"."+field+" = ?", true)
			} else {
				db = db.Where(modelType+"."+field+" = ?", false)
			}
		} else if isTimeField(field) {
			db = db.Where(modelType+"."+field+" <= ?", field)
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

func isBooleanField(field string) bool {
	switch field {
	case "is_completed":
		return true
	default:
		return false
	}
}

func isTimeField(field string) bool {
	switch field {
	case "created_at", "deadline":
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
