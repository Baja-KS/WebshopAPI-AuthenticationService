package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	seedTmp:=os.Getenv("SEED_IF_EMPTY")

	seed:=false

	if seedTmp=="true" {
		seed=true
	}

	if !seed {
		return nil
	}
	var count int64
	count=0
	db.Model(&User{}).Count(&count)

	if count!=0 {
		return nil
	}

	hash,_:=HashPassword("testadmin")
	admin:=User{
		Username:  "testadmin",
		Fullname:  "Test Admin",
		Email:     "testadmin@example.com",
		Password:  hash,
		IsAdmin:   true,
	}
	db.Create(&admin)
	hash,_=HashPassword("testuser")
	user:=User{
		Username:  "testuser",
		Fullname:  "Test User",
		Email:     "testuser@example.com",
		Password:  hash,
		IsAdmin:   false,
	}
	db.Create(&user)

	return nil
}
