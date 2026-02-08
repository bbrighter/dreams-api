package migrations

import (
	"strconv"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration11 = &gormigrate.Migration{
	ID: "20260129_ensure_unqiueness",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct{ ID uint }
		type CategoryType string
		type Category struct {
			ID     uint
			Name   string       `gorm:"uniqueIndex:idx_category_type"`
			Type   CategoryType `gorm:"uniqueIndex:idx_category_type"`
			Dreams []Dream      `gorm:"many2many:categories_dreams"`
		}

		type result struct {
			Name  string
			Count int
		}
		var results []result
		if err := tx.Model(&Category{}).
			Select("name", "count(*) as count").
			Group("name").Having("count > ?", 1).
			Find(&results).Error; err != nil {
			return err
		}
		for _, res := range results {
			var cats []Category
			if err := tx.Where("name = ?", res.Name).Find(&cats).Error; err != nil {
				return err
			}
			for i, cat := range cats {
				if err := tx.Model(&Category{}).Where("id = ?", cat.ID).UpdateColumn("name", cat.Name+"_"+strconv.Itoa(i)).Error; err != nil {
					return err
				}
			}
		}
		return tx.AutoMigrate(&Category{})

	},
	Rollback: func(tx *gorm.DB) error {
		type Category struct{ ID uint }
		return tx.Migrator().DropIndex(&Category{}, "idx_category_type")
	},
}
