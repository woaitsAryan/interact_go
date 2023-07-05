package utils

import "gorm.io/gorm"

func PostSelectConfig(db *gorm.DB) *gorm.DB {
	return db.Select("id, content, created_at, tags, images, user_id, no_likes, edited")
}
