package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Search(c *fiber.Ctx, index int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchStr := c.Query("search", "")
		if searchStr == "" {
			return db
		}

		searchStr = strings.ToLower(searchStr)

		substrings := strings.Split(searchStr, " ")
		var regexPatterns []string
		for _, substring := range substrings {
			pattern := ".*" + substring + ".*"
			regexPatterns = append(regexPatterns, pattern)
		}

		switch index {
		case 0: //* users
			// for _, pattern := range regexPatterns {
			// 	db = db.Or("LOWER(name) ~ ? OR LOWER(username) ~ ? OR ? = ANY (tags)", pattern, pattern, searchStr)
			// }
			db = db.Where("LOWER(name) LIKE ? OR LOWER(username) LIKE ? OR ? = ANY (tags)", "%"+searchStr+"%", "%"+searchStr+"%", searchStr)
			return db
		case 1: //* projects
			// for _, pattern := range regexPatterns {
			// 	db = db.Or("LOWER(title) ~ ? OR ? = ANY (tags)", pattern, searchStr)
			// }
			db = db.Where("LOWER(title) LIKE ? OR ? = ANY (tags)", "%"+searchStr+"%", searchStr)
			return db
		case 2: //* posts
			for _, pattern := range regexPatterns {
				db = db.Or("LOWER(content) ~ ? OR ? = ANY (tags) ", pattern, searchStr)
			}
			return db
		case 3: //* openings
			db = db.Joins("JOIN projects ON openings.project_id = projects.id").
				Where("LOWER(openings.title) LIKE ? OR LOWER(projects.title) LIKE ? OR ? = ANY (openings.tags) OR ? = ANY (projects.tags)",
					"%"+searchStr+"%", "%"+searchStr+"%", searchStr, searchStr)
			return db
		case 4: //* search_queries
			db = db.Where("LOWER(query) LIKE ?", "%"+searchStr+"%")
			return db
		default:
			return db
		}
	}
}
