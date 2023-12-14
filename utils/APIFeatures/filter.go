package utils

import (
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func genericFilter(db *gorm.DB, field, value string) *gorm.DB {
	if value != "" {
		if isUserField(field) {
			db = db.Joins("User").Where("user."+field+" ILIKE ?", "%"+value+"%")
		} else if isArrayField(field) {
			db = db.Where("? ILIKE ANY ("+field+")", "%"+value+"%")
		} else if isOrg(field) {
			db = db.Joins("Organization").Where("organization.organization_title ILIKE ?", "%"+value+"%")
		} else {
			db = db.Where(field+" ILIKE ?", "%"+value+"%")
		}
	}
	return db
}

func eventTimeSearch(db *gorm.DB, eventTime string) *gorm.DB {
	if eventTime != "" {
		parsedEventTime, err := time.Parse(time.RFC3339, eventTime)
		if err != nil {
			helpers.LogServerError("Error parsing start timestamp", err, "timestampSearch")
			return db
		}
		return db.Where("StartTime <= ? AND EndTime >= ?", parsedEventTime, parsedEventTime)
	}
	return db

}

func timestampSearch(db *gorm.DB, start, end string) *gorm.DB {
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

		return db.Where("CreatedAt BETWEEN ? AND ?", startTime, endTime)
	} else if start != "" {
		startTime, err := time.Parse(time.RFC3339, start)
		if err != nil {
			helpers.LogServerError("Error parsing start timestamp", err, "timestampSearch")
			return db
		}

		return db.Where("CreatedAt >= ?", startTime)
	} else if end != "" {
		endTime, err := time.Parse(time.RFC3339, end)
		if err != nil {
			helpers.LogServerError("Error parsing end timestamp", err, "timestampSearch")
			return db
		}

		return db.Where("CreatedAt <= ?", endTime)
	}
	return db
}

func Filter(c *fiber.Ctx, index int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch index {
		//* For Project
		case 1:
			projectFields := []string{"title", "category", "tags", "username", "user"}
			for _, projectField := range projectFields {
				value := c.Query(projectField, "")
				db = genericFilter(db, projectField, value)
			}

			startTime := c.Query("start", "")
			endTime := c.Query("end", "")
			db = timestampSearch(db, startTime, endTime)
			return db

		//* For Post
		case 2:
			postFields := []string{"tags", "username", "user"}

			for _, postField := range postFields {
				value := c.Query(postField, "")
				db = genericFilter(db, postField, value)
			}

			startTime := c.Query("start", "")
			endTime := c.Query("end", "")
			db = timestampSearch(db, startTime, endTime)
			return db

		//* For Event
		case 3:
			eventFields := []string{"title", "tagline", "org", "tags"}

			for _, eventField := range eventFields {
				value := c.Query(eventField, "")
				db = genericFilter(db, eventField, value)
			}
			eventTime := c.Query("eventTime", "")
			db = eventTimeSearch(db, eventTime)
			return db

		//* For opening
		case 4:
			openingFields := []string{"title", "tags", "username", "user"}

			for _, openingField := range openingFields {
				value := c.Query(openingField, "")
				db = genericFilter(db, openingField, value)
			}

			startTime := c.Query("start", "")
			endTime := c.Query("end", "")
			db = timestampSearch(db, startTime, endTime)
			return db
		default:
			return db
		}
	}
}

func isArrayField(field string) bool {
	switch field {
	case "tags":
		return true
	default:
		return false
	}
}

func isUserField(field string) bool {
	switch field {
	case "username", "name":
		return true
	default:
		return false
	}
}

func isOrg(field string) bool {
	switch field {
	case "org":
		return true
	default:
		return false
	}
}
