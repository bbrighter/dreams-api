package entities

type CategoryType string

const (
	TypePerson   CategoryType = "person"
	TypeCategory CategoryType = "category"
)

type Category struct {
	ID     uint
	Name   string       `gorm:"uniqueIndex:idx_category_type"`
	Type   CategoryType `gorm:"uniqueIndex:idx_category_type"`
	Dreams []Dream      `gorm:"many2many:categories_dreams"`
}

type Categories []Category
