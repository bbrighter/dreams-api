package store

import (
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const (
	initialMigration = "20231218_Initial"
	migration1       = "20240105_Tags"
	migration2       = "20240131_Persons"
	migration3       = "20240207_RenameTags"
	migration4       = "20240218_HiddenDreams"
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
				Dreams []Dream `gorm:"many2many:tags_dreams;"`
			}
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Tags        []Tag `gorm:"many2many:tags_dreams;"`
			}
			if err := tx.Migrator().DropTable(Tag{}); err != nil {
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
				Dreams []Dream `gorm:"many2many:tags_dreams;"`
			}
			type Person struct {
				ID     uint
				Name   string
				Dreams []Dream `gorm:"many2many:people_dreams;"`
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
			type Tag struct {
				ID     uint
				Title  string
				Dreams []Dream `gorm:"many2many:tags_dreams;"`
			}
			if err := tx.Migrator().RenameColumn(&Tag{}, "title", "name"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("tags", "categories"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("tags_dreams", "categories_dreams"); err != nil {
				return err
			}
			type CategoriesDream struct {
				TagId uint
			}
			if err := tx.Migrator().RenameColumn(&CategoriesDream{}, "tag_id", "category_id"); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			type Category struct {
				ID     uint
				Name   string
				Dreams []Dream `gorm:"many2many:categories_dreams;"`
			}
			if err := tx.Migrator().RenameColumn(&Category{}, "name", "title"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("categories", "tags"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("categories_dreams", "tags_dreams"); err != nil {
				return err
			}
			type TagsDream struct {
				TagId uint
			}
			if err := tx.Migrator().RenameColumn(&TagsDream{}, "category_id", "tag_id"); err != nil {
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
				Visible     *bool      `gorm:"default:true"`
				Categories  []Category `gorm:"many2many:categories_dreams;"`
				Persons     []Person   `gorm:"many2many:people_dreams;"`
			}
			return tx.Migrator().AutoMigrate(&Dream{})
		},
		Rollback: func(tx *gorm.DB) error {
			type Dream struct {
				ID          uint
				Date        time.Time
				Description string
				Visible     *bool      `gorm:"default:true"`
				Categories  []Category `gorm:"many2many:categories_dreams;"`
				Persons     []Person   `gorm:"many2many:people_dreams;"`
			}
			return tx.Migrator().DropColumn(&Dream{}, "visible")
		},
	},
}

func migrationFactory(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}

func Migration(repo Repo) error {
	m := migrationFactory(repo.db)
	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate")
		return err
	}
	log.Printf("Migration succesful")
	return nil
}

func RollbackTo(db *gorm.DB, rollbackTo string) error {
	m := migrationFactory(db)
	if err := m.RollbackTo(rollbackTo); err != nil {
		log.Fatalf("Could not rollback")
		return err
	}
	println("Rollback to " + rollbackTo)
	return nil
}

func Rollback(repo Repo) error {
	m := migrationFactory(repo.db)
	if err := m.RollbackLast(); err != nil {
		log.Fatalf("Could not rollback")
		return err
	}
	println("Rollback succesful")
	return nil
}
