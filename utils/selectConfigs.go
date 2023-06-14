package utils

import "gorm.io/gorm"

func PostSelectConfig(db *gorm.DB) *gorm.DB {
	return db.Select("id, content, posted_at, tags, images, users.id, users.username, users.name, users.profile_pic")
}
