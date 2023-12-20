package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func InitRepo(dbName string) Repo {
	db := initDatabase(dbName)
	return Repo{db: db}
}

func initDatabase(dbName string) *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	return db
}
