package repository

import (
	"os"
	"path"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func newDatabase(name string, log *zap.Logger) *gorm.DB {
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
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot open database", zap.Error(err))
	}
	return db
}
