package utils

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type APIFeatures struct {
	Context *fiber.Ctx
	DB      *gorm.DB
	Model   interface{}
}

func NewAPIFeatures(ctx *fiber.Ctx, db *gorm.DB, model interface{}) *APIFeatures {
	return &APIFeatures{
		Context: ctx,
		DB:      db,
		Model:   model,
	}
}

func (af *APIFeatures) Search() *APIFeatures {
	searchStr := af.Context.Query("search")

	if searchStr != "" {
		// Implement your search logic using GORM's Where clause
		af.DB = af.DB.Where("column LIKE ?", "%"+searchStr+"%")
	}

	return af
}

func (af *APIFeatures) Filter() *APIFeatures {
	// Implement your filtering logic here using GORM's Where clause
	return af
}

func (af *APIFeatures) Sort() *APIFeatures {
	sortStr := af.Context.Query("sort")
	if sortStr != "" {
		// Implement your sorting logic here using GORM's Order clause
		af.DB = af.DB.Order(sortStr)
	}

	return af
}

func (af *APIFeatures) Fields() *APIFeatures {
	// Implement your field selection logic here using GORM's Select clause
	return af
}

func (af *APIFeatures) Paginator() *APIFeatures {
	// page := af.Context.Query("page")
	// limit := af.Context.Query("limit")

	// Implement your pagination logic here using GORM's Offset and Limit clauses
	return af
}
