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

		var searchCondition interface{}
		switch index {
		case 0: // !users
			searchCondition = gorm.Expr("$or", []interface{}{
				map[string]interface{}{
					"username": gorm.Expr("$in", regexArry),
				},
				map[string]interface{}{
					"name": gorm.Expr("$in", regexArry),
				},
			})
		case 1: // !projects
			searchCondition = gorm.Expr("$or", []interface{}{
				map[string]interface{}{
					"title": gorm.Expr("$in", regexArry),
				},
				map[string]interface{}{
					"tags": gorm.Expr("$elemMatch", map[string]interface{}{
						"$in": regexArry,
					}),
				},
				map[string]interface{}{
					"category": gorm.Expr("$in", regexArry),
				},
			})
		case 2: // !posts
			searchCondition = gorm.Expr("$or", []interface{}{
				map[string]interface{}{
					"content": gorm.Expr("$in", regexArry),
				},
				map[string]interface{}{
					"tags": gorm.Expr("$elemMatch", map[string]interface{}{
						"$in": regexArry,
					}),
				},
			})
		case 3: // !openings
			searchCondition = gorm.Expr("$or", []interface{}{
				map[string]interface{}{
					"title": gorm.Expr("$in", regexArry),
				},
				map[string]interface{}{
					"description": gorm.Expr("$in", regexArry),
				},
				map[string]interface{}{
					"tags": gorm.Expr("$elemMatch", map[string]interface{}{
						"$in": regexArry,
					}),
				},
			})
		default:
			searchCondition = nil
		}

		if searchCondition != nil {
			db = db.Where(searchCondition)
		}

		return db
	}
}
