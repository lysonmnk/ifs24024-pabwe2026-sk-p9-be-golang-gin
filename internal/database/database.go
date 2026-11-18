package database

import (
	"fmt"

	"app/internal/config"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var AutoMigrator = AutoMigrate
var AdminSeeder = SeedAdmin

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	if cfg.DBConnection == "sqlite" {
		dialector = sqlite.Open(cfg.DBDatabase)
	} else {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUsername, cfg.DBPassword, cfg.DBDatabase, cfg.DBPort, cfg.DBSSLMode,
		)
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := AutoMigrator(db); err != nil {
		return nil, err
	}

	if err := AdminSeeder(db, cfg); err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Post{})
}

func SeedAdmin(db *gorm.DB, cfg *config.Config) error {
	var count int64
	if err := db.Model(&User{}).Where("username = ? OR role = ?", cfg.AdminUsername, "admin").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := User{
		Name:     "Administrator",
		Username: cfg.AdminUsername,
		Email:    cfg.AdminUsername + "@delcom.org",
		Password: string(hashedPassword),
		Role:     "admin",
	}
	return db.Create(&admin).Error
}
