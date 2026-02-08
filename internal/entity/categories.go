package entity

import (
	"gorm.io/gorm"
)

type CategoryType string

const (
	TypePerson   CategoryType = "person"
	TypeCategory CategoryType = "category"
)

var validCategoryTypes = map[CategoryType]bool{
	TypePerson:   true,
	TypeCategory: true,
}

func NewCategoryType(str string) (CategoryType, error) {
	if !validCategoryTypes[CategoryType(str)] {
		return TypePerson, ErrorBadParamWithReasons("invalid category type")
	}
	return CategoryType(str), nil
}

type Category struct {
	ID     uint
	Name   string       `gorm:"uniqueIndex:idx_category_type"`
	Type   CategoryType `gorm:"uniqueIndex:idx_category_type"`
	Dreams []Dream      `gorm:"many2many:categories_dreams"`
}

type Categories []Category

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories" binding:"required"`
}

type CategoryResponse struct {
	ID   uint         `json:"id" binding:"required"`
	Name string       `json:"name" binding:"required"`
	Type CategoryType `json:"type" binding:"required"`
}

func (c Category) ToResponse() CategoryResponse {
	return CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
		Type: c.Type,
	}
}
func (cats Categories) ToResponse() CategoriesResponse {
	var categories []CategoryResponse
	for _, c := range cats {
		categories = append(categories, c.ToResponse())
	}
	return CategoriesResponse{Categories: categories}
}

func (cats Categories) ToList(t CategoryType) []CategoryResponse {
	var list []CategoryResponse
	for _, c := range cats {
		list = append(list, c.ToResponse())
	}
	return list
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if !validCategoryTypes[c.Type] {
		return ErrorBadParamWithReasons("invalid category type")
	}
	return nil
}

func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	if !validCategoryTypes[c.Type] {
		return ErrorBadParamWithReasons("invalid category type")
	}
	return nil
}
