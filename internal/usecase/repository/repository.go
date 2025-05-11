package repository

import (
	"os"
	"path"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func NewDatabase(name string, log *zap.Logger) *gorm.DB {
	var dbName string
	if name == ":memory:" {
		dbName = name
	} else {
		dbName = path.Join("data", name)
		err := os.MkdirAll(path.Join(".", "data"), os.ModePerm)
		if err != nil {
			log.Fatal("Cannot create database folder", zap.Error(err))
		}
	}
	logger := zapgorm2.New(log)
	logger.SetAsDefault()
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{Logger: logger})
	if err != nil {
		log.Fatal("Cannot open database", zap.Error(err))
	}
	return db
}
