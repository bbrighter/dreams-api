package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration13 = &gormigrate.Migration{
	ID: "20261016_addTranscripts",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct {
			ID uint
		}
		type RecordingState string
		type TranscriptionState string
		type Recording struct {
			ID                 uint
			Dream              Dream `gorm:"constraint:OnDelete:CASCADE;"`
			DreamID            uint  `gorm:"unique"`
			RecordingState     RecordingState
			TranscriptionState TranscriptionState
			Transcript         string
			MimeType           string
		}

		return tx.AutoMigrate(&Recording{})

	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("recording")
	},
}
