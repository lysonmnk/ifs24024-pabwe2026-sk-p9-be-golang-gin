package database

import (
	"fmt"
	"time"

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

	if cfg.DBConnection == "sqlite" || cfg.DBConnection == "none" || cfg.DBHost == "delcom-deployment-postgres" {
		dbPath := cfg.DBDatabase
		if dbPath == "" || dbPath == "pabwe_p9_blog" || cfg.DBHost == "delcom-deployment-postgres" {
			dbPath = "pabwe_p9_blog.db"
		}
		dialector = sqlite.Open(dbPath)
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
	_ = db.Model(&User{}).Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR username = ?", "admin%", "admin%", cfg.AdminUsername).Update("role", "admin")

	var adminUser User
	if err := db.Where("username = ?", cfg.AdminUsername).First(&adminUser).Error; err == nil {
		if adminUser.Role != "admin" {
			adminUser.Role = "admin"
			_ = db.Save(&adminUser)
		}
	} else {
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
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	}

	var postCount int64
	if err := db.Model(&Post{}).Count(&postCount).Error; err == nil && postCount == 0 {
		var adminUser User
		if err := db.Where("username = ? OR role = ?", cfg.AdminUsername, "admin").First(&adminUser).Error; err == nil {
			now := time.Now()
			initialPost := Post{
				UserID:      adminUser.ID,
				Title:       "Pengenalan Pemrograman Golang untuk Backend Modern",
				Slug:        "pengenalan-pemrograman-golang-untuk-backend-modern",
				Category:    "Teknologi",
				Excerpt:     "Panduan dasar pengembangan aplikasi web berperforma tinggi menggunakan Golang.",
				Content:     "Golang adalah bahasa pemrograman yang dirancang untuk efisiensi, konkurensi, dan performa tinggi pada sistem backend modern.",
				ReadingTime: 2,
				IsPublished: true,
				PublishedAt: &now,
			}
			_ = db.Create(&initialPost)
		}
	}

	return nil
}
