package entities

import (
	"time"
)

type Dream struct {
	ID          uint
	Date        time.Time
	Description string
	Finalized   bool
	Categories  Categories `gorm:"many2many:categories_dreams"`
	Rating      *int
}

type Dreams []Dream
