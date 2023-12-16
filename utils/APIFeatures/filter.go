package utils

import (
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func timestampSearch(db *gorm.DB, start, end string, modelType string) *gorm.DB {
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

		return db.Where(modelType + ".created_at BETWEEN ? AND ?", startTime, endTime)
	} else if start != "" {
		startTime, err := time.Parse(time.RFC3339, start)
		if err != nil {
			helpers.LogServerError("Error parsing start timestamp", err, "timestampSearch")
			return db
		}

		return db.Where(modelType +".created_at >= ?", startTime)
	} else if end != "" {
		endTime, err := time.Parse(time.RFC3339, end)
		if err != nil {
			helpers.LogServerError("Error parsing end timestamp", err, "timestampSearch")
			return db
		}

		return db.Where(modelType + ".created_at <= ?", endTime)
	}
	return db
}

func Filter(c *fiber.Ctx, index int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch index {
		//* For Project
		case 1:
			projectFields := []string{"title", "category", "tags", "username", "name"}
			for _, projectField := range projectFields {
				value := c.Query(projectField, "")
				db = genericFilter(db, projectField, value, "projects")
			}

			startTime := c.Query("start", "")
			endTime := c.Query("end", "")
			db = timestampSearch(db, startTime, endTime, "projects")
			return db

		//* For User
		case 2:
			userFields := []string{"tags", "location", "areas_of_collaboration", "school"}
			// tags is skills, roles is areasOfCollaboration, school is college
			for _, userField := range userFields {
				value := c.Query(userField, "")
				db = genericFilter(db, userField, value, "users")
			}

			startTime := c.Query("start", "")
			endTime := c.Query("end", "")
			db = timestampSearch(db, startTime, endTime, "users")
			return db
 
		//* For Event
		case 3:
			eventFields := []string{"title", "tagline", "org", "tags", "location"}

			for _, eventField := range eventFields {
				value := c.Query(eventField, "")
				db = genericFilter(db, eventField, value, "events")
			}
			eventTime := c.Query("eventTime", "")
			db = eventTimeSearch(db, eventTime)
			return db

		//* For opening
		case 4:
			openingFields := []string{"title", "tags"}

			for _, openingField := range openingFields {
				value := c.Query(openingField, "")
				db = genericFilter(db, openingField, value, "openings")
			}

			startTime := c.Query("start", "")
			endTime := c.Query("end", "")
			db = timestampSearch(db, startTime, endTime, "openings")
			return db
		default:
			return db
		}
	}
}


func eventTimeSearch(db *gorm.DB, eventTime string) *gorm.DB {
	if eventTime != "" {
		parsedEventTime, err := time.Parse(time.RFC3339, eventTime)
		if err != nil {
			helpers.LogServerError("Error parsing start timestamp", err, "timestampSearch")
			return db
		}
		return db.Where("start_time <= ? AND end_time >= ?", parsedEventTime, parsedEventTime)
	}
	return db
}

func genericFilter(db *gorm.DB, field, value string, modelType string) *gorm.DB {
	if value != "" {
		if isUserField(field) {
			db = db.Joins("JOIN users ON " + modelType + ".user_id = users.id" ).Where("users." + field +" ILIKE ?", "%"+value+"%")
		} else if isArrayField(field) {
			db = db.Where("? =  ANY("+modelType+"."+field+")",value)
		} else if isOrg(field) {
			db = db.Joins("JOIN organizations ON " + modelType + ".organization_id = organizations.id").Where("organizations.organization_title ILIKE ?", "%"+value+"%")
		} else if isProfile(field) {
			if(field == "areas_of_collaboration") {
				db = db.Joins("JOIN profiles ON users.id = profiles.user_id").Where("? =  ANY("+"profiles."+field+")",value)
			} else{
				db = db.Joins("JOIN profiles ON users.id = profiles.user_id").Where("profiles." + field +" ILIKE ?", "%"+value+"%")
			}
		} else {
			db = db.Where(modelType+"." +field+" ILIKE ?", "%"+value+"%")
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

func isProfile(field string) bool {
	switch field {
	case "location", "areas_of_collaboration", "school":
		return true
	default:
		return false
	}
}