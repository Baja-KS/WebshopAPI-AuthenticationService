package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)
import log "github.com/sirupsen/logrus"

func NewDatabase() (*gorm.DB, error) {
	log.Info("Setting up database connection")
	dbUsername:=os.Getenv("DB_USERNAME")
	dbPassword:=os.Getenv("DB_PASSWORD")
	dbHost:=os.Getenv("DB_HOST")
	dbName:=os.Getenv("DB_NAME")
	dbPort:=os.Getenv("DB_PORT")
	sslMode:=os.Getenv("SSL_MODE")

	connectString :=fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",dbHost,dbPort,dbUsername,dbName,dbPassword,sslMode)

	db,err:=gorm.Open(postgres.Open(connectString),&gorm.Config{})
	//if err != nil {
	//	return db, err
	//}
	for err!=nil {
		log.Info("Reconnecting to database...")
		db,err=gorm.Open(postgres.Open(connectString),&gorm.Config{})
	}

	if sqlDB,err:=db.DB();err!=nil || sqlDB.Ping()!=nil{
		return db,err
	}

	return db, nil
}
