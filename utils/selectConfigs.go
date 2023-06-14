package utils

import "gorm.io/gorm"

func PostSelectConfig(db *gorm.DB) *gorm.DB {
	return db.Select("id, content, posted_at, tags, images, user_id")
}
