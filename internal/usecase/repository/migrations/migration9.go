package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration9 = &gormigrate.Migration{
	ID: "20250612_SimplifyCategories",
	Migrate: func(tx *gorm.DB) error {
		type Category struct {
			ID   uint
			Name string
			Type string
		}
		if !tx.Migrator().HasColumn(&Category{}, "Type") {
			if err := tx.Migrator().AddColumn(&Category{}, "Type"); err != nil {
				return err
			}
		}

		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&Category{}).Update("Type", "category").Error; err != nil {
			return err
		}

		type Person struct {
			ID   uint
			Name string
		}

		var persons []Person
		if err := tx.Find(&persons).Error; err != nil {
			return err
		}
		if len(persons) > 0 {
			var cats []Category
			var oldIds = make(map[uint]string)
			for _, p := range persons {
				cat := Category{Name: p.Name, Type: "person"}
				oldIds[p.ID] = p.Name
				cats = append(cats, cat)
			}
			if err := tx.Create(&cats).Error; err != nil {
				return err
			}

			type PeopleDream struct {
				PersonID uint
				DreamID  uint
			}
			var peopleDream []PeopleDream
			if err := tx.Find(&peopleDream).Error; err != nil {
				return err
			}
			type CategoriesDream struct {
				CategoryID uint
				DreamID    uint
			}
			var categoriesDream []CategoriesDream
			for _, pd := range peopleDream {
				for _, cat := range cats {
					if cat.Name == oldIds[pd.PersonID] {
						cd := CategoriesDream{DreamID: pd.DreamID, CategoryID: cat.ID}
						categoriesDream = append(categoriesDream, cd)
					}
				}
			}
			if err := tx.Create(&categoriesDream).Error; err != nil {
				return err
			}
		}

		type PeopleDream struct {
			PersonID uint
			DreamID  uint
		}
		if err := tx.Migrator().DropTable(&PeopleDream{}); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable(&Person{}); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(d *gorm.DB) error {
		return nil
	},
}
