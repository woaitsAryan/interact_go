package utils

import (
	"regexp"
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

		regexArry := make([]string, 0)
		searchComponents := strings.Split(searchStr, " ")

		for _, item := range searchComponents {
			regexArry = append(regexArry, regexp.QuoteMeta(item))
		}

		regexArry = append(regexArry, regexp.QuoteMeta(searchStr))
		regexArry = append(regexArry, regexp.QuoteMeta(strings.ReplaceAll(searchStr, " ", "")))

		interfaceArry := make([]interface{}, len(regexArry))
		for i, v := range regexArry {
			interfaceArry[i] = v
		}

		searchStr = strings.ToLower(searchStr)

		// var searchCondition interface{}
		switch index {
		case 0: //* users
			// searchCondition = []interface{}{
			// 	map[string]interface{}{
			// 		"username": gorm.Expr("IN (?)", interfaceArry),
			// 	},
			// 	map[string]interface{}{
			// 		"name": gorm.Expr("IN (?)", interfaceArry),
			// 	},
			// }
			db = db.Where("LOWER(name) LIKE ? OR LOWER(username) LIKE ? OR ? = ANY (tags)", "%"+searchStr+"%", "%"+searchStr+"%", searchStr)
			return db
		case 1: //* projects
			// searchCondition = []interface{}{
			// 	map[string]interface{}{
			// 		"title": gorm.Expr("IN (?)", interfaceArry),
			// 	},
			// 	map[string]interface{}{
			// 		"tags": gorm.Expr("$elemMatch", map[string]interface{}{
			// 			"$in": interfaceArry,
			// 		}),
			// 	},
			// 	map[string]interface{}{
			// 		"category": gorm.Expr("IN (?)", interfaceArry),
			// 	},
			// }
			db = db.Where("LOWER(title) LIKE ? OR ? = ANY (tags)", "%"+searchStr+"%", searchStr)
			return db
		case 2: //* posts
			// searchCondition = []interface{}{
			// 	map[string]interface{}{
			// 		"content": gorm.Expr("content IN (?)", interfaceArry),
			// 	},
			// 	map[string]interface{}{
			// 		"tags": gorm.Expr("$elemMatch", map[string]interface{}{
			// 			"$in": interfaceArry,
			// 		}),
			// 	},
			// }
			db = db.Where("LOWER(content) LIKE ? OR ? = ANY (tags) ", "%"+searchStr+"%", searchStr)
			return db
		case 3: //* openings
			// searchCondition = []interface{}{
			// 	map[string]interface{}{
			// 		"title": gorm.Expr("IN (?)", interfaceArry),
			// 	},
			// 	map[string]interface{}{
			// 		"description": gorm.Expr("IN (?)", interfaceArry),
			// 	},
			// 	map[string]interface{}{
			// 		"tags": gorm.Expr("$elemMatch", map[string]interface{}{
			// 			"$in": interfaceArry,
			// 		}),
			// 	},
			// }
			db = db.Where("LOWER(title) LIKE ? OR ? = ANY (tags)  ", "%"+searchStr+"%", searchStr)
			return db
		default:
			return db
			// searchCondition = nil
		}
	}
}
