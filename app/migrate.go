package app

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
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     *bool             `gorm:"default:true"`
				Categories  []entity.Category `gorm:"many2many:categories_dreams;"`
				Persons     []entity.Person   `gorm:"many2many:people_dreams;"`
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
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     bool              `gorm:"default:true"`
				Categories  []entity.Category `gorm:"many2many:categories_dreams;"`
				Persons     []entity.Person   `gorm:"many2many:people_dreams;"`
			}
			return tx.AutoMigrate(&Dream{})
		},
		Rollback: func(tx *gorm.DB) error {
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     *bool             `gorm:"default:true"`
				Categories  []entity.Category `gorm:"many2many:categories_dreams;"`
				Persons     []entity.Person   `gorm:"many2many:people_dreams;"`
			}
			return tx.AutoMigrate(&Dream{})
		},
	},
}

func migrationFactory(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}

func migration(db *gorm.DB, logger *zap.Logger) {
	m := migrationFactory(db)
	if err := m.Migrate(); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}
	logger.Info("Migration successful")
}
