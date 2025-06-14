package repository

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	initialMigration = "20231218_Initial"
	migration1       = "20240105_Tags"
	migration2       = "20240131_Persons"
	migration3       = "20240207_RenameTags"
	migration4       = "20240218_HiddenDreams"
	migration5       = "20240728_VisibleNonPointer"
	migration6       = "20250601_FinalizeDreams"
	migration7       = "20250603_FixFinalizeDreams"
	migration8       = "20250612_SimplifyCategories"
)

var migrations = []*gormigrate.Migration{
	{
		ID: initialMigration,
		Migrate: func(tx *gorm.DB) error {
			type Tag struct {
				ID      uint
				Title   string
				DreamID uint
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Tags        []*Tag
			}

			return tx.AutoMigrate(
				Dream{},
				Tag{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable("dreams"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("tags"); err != nil {
				return err
			}
			return nil
		},
	},
	{
		ID: migration1,
		Migrate: func(tx *gorm.DB) error {
			type Tag struct {
				ID     uint
				Title  string
				Dreams []entity.Dream `gorm:"many2many:tags_dreams;"`
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Tags        []Tag `gorm:"many2many:tags_dreams;"`
			}
			if err := tx.Migrator().DropTable("tags"); err != nil {
				return err
			}
			return tx.AutoMigrate(
				Dream{},
				Tag{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			type Tag struct {
				ID      uint
				Title   string
				DreamID uint
			}
			if err := tx.Migrator().DropTable("tags"); err != nil {
				return err
			}
			if err := tx.AutoMigrate(Tag{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("tags_dreams"); err != nil {
				return err
			}
			return nil
		},
	},
	{
		ID: migration2,
		Migrate: func(tx *gorm.DB) error {
			type Tag struct {
				ID     uint
				Title  string
				Dreams []entity.Dream `gorm:"many2many:tags_dreams;"`
			}
			type Person struct {
				ID     uint
				Name   string
				Dreams []entity.Dream `gorm:"many2many:people_dreams;"`
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Tags        []Tag    `gorm:"many2many:tags_dreams;"`
				Persons     []Person `gorm:"many2many:people_dreams;"`
			}
			return tx.Migrator().AutoMigrate(Person{}, Dream{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable("people_dreams"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("people"); err != nil {
				return err
			}
			return nil
		},
	},
	{
		ID: migration3,
		Migrate: func(tx *gorm.DB) error {
			type Tag struct{}
			if err := tx.Migrator().RenameColumn(Tag{}, "title", "name"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("tags", "categories"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("tags_dreams", "categories_dreams"); err != nil {
				return err
			}
			type CategoriesDreams struct{}
			if err := tx.Migrator().RenameColumn(CategoriesDreams{}, "tag_id", "category_id"); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().RenameColumn("categories", "name", "title"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("categories", "tags"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("categories_dreams", "tags_dreams"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameColumn("tags_dreams", "category_id", "tag_id"); err != nil {
				return err
			}
			return nil
		},
	},
	{
		ID: migration4,
		Migrate: func(tx *gorm.DB) error {
			type Person struct {
				ID   uint
				Name string
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     *bool             `gorm:"default:true"`
				Categories  []entity.Category `gorm:"many2many:categories_dreams;"`
				Persons     []Person          `gorm:"many2many:people_dreams;"`
			}
			return tx.Migrator().AutoMigrate(&Dream{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn("dreams", "visible")
		},
	},
	{
		ID: migration5,
		Migrate: func(tx *gorm.DB) error {
			type Person struct {
				ID   uint
				Name string
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     bool              `gorm:"default:true"`
				Categories  []entity.Category `gorm:"many2many:categories_dreams;"`
				Persons     []Person          `gorm:"many2many:people_dreams;"`
			}
			return tx.AutoMigrate(&Dream{})
		},
		Rollback: func(tx *gorm.DB) error {
			type Person struct {
				ID   uint
				Name string
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     *bool             `gorm:"default:true"`
				Categories  []entity.Category `gorm:"many2many:categories_dreams;"`
				Persons     []Person          `gorm:"many2many:people_dreams;"`
			}
			return tx.AutoMigrate(&Dream{})
		},
	},
	{
		ID: migration6,
		Migrate: func(tx *gorm.DB) error {
			type Person struct {
				ID   uint
				Name string
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     bool `gorm:"default:true"`
				Finalized   bool
				Categories  []entity.Category `gorm:"many2many:categories_dreams;"`
				Persons     []Person          `gorm:"many2many:people_dreams;"`
			}
			if err := tx.AutoMigrate(&Dream{}); err != nil {
				return err
			}
			var dreams []Dream
			tx.Preload("Categories").Preload("Persons").
				Where(`
					dreams.id NOT IN (
						SELECT dream_id FROM categories_dreams
					)
					OR dreams.id NOT IN (
						SELECT dream_id FROM people_dreams
					)
				`).Find(&dreams)
			var ids []uint
			for _, d := range dreams {
				ids = append(ids, d.ID)
			}
			return tx.Model(&Dream{}).Where("id IN ?", ids).UpdateColumn("finalized", true).Error
		},
		Rollback: func(tx *gorm.DB) error {
			type Person struct {
				ID   uint
				Name string
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     bool              `gorm:"default:true"`
				Categories  []entity.Category `gorm:"many2many:categories_dreams;"`
				Persons     []Person          `gorm:"many2many:people_dreams;"`
			}
			return tx.AutoMigrate(&Dream{})
		},
	},
	{
		ID: migration7,
		Migrate: func(tx *gorm.DB) error {
			type Dream struct{}
			return tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&Dream{}).Update("finalized", gorm.Expr("NOT finalized")).Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	},
	{
		ID: migration8,
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
	},
}

func migrationFactory(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}

func Migration(db *gorm.DB, logger *zap.Logger) {
	m := migrationFactory(db)
	if err := m.Migrate(); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}
	logger.Info("Migration successful")
}
