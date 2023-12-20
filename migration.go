package main

import (
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const (
	initialMigration = "20231218_Initial"
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
}

func migrationFactory(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}

func Migration(db *gorm.DB) error {
	m := migrationFactory(db)
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

func Rollback(db *gorm.DB) error {
	m := migrationFactory(db)
	if err := m.RollbackLast(); err != nil {
		log.Fatalf("Could not rollback")
		return err
	}
	println("Rollback succesful")
	return nil
}
