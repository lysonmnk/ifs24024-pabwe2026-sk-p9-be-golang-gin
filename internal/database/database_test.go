package database

import (
	"errors"
	"strings"
	"testing"

	"app/internal/config"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestDatabase(t *testing.T) {
	cfg := &config.Config{
		DBConnection:  "sqlite",
		DBDatabase:    ":memory:",
		AdminUsername: "admin",
		AdminPassword: "adminpassword123",
	}

	t.Run("InitDB success with sqlite", func(t *testing.T) {
		db, err := InitDB(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, db)

		var admin User
		err = db.Where("username = ?", "admin").First(&admin).Error
		assert.NoError(t, err)
		assert.Equal(t, "admin", admin.Role)
	})

	t.Run("InitDB fails when AutoMigrator returns error", func(t *testing.T) {
		orig := AutoMigrator
		defer func() { AutoMigrator = orig }()

		AutoMigrator = func(db *gorm.DB) error {
			return errors.New("migration failed")
		}

		db, err := InitDB(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("InitDB fails when AdminSeeder returns error", func(t *testing.T) {
		longPassCfg := &config.Config{
			DBConnection:  "sqlite",
			DBDatabase:    ":memory:",
			AdminUsername: "admin",
			AdminPassword: strings.Repeat("a", 100),
		}

		db, err := InitDB(longPassCfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("SeedAdmin when admin already exists", func(t *testing.T) {
		db, err := InitDB(cfg)
		assert.NoError(t, err)

		err = SeedAdmin(db, cfg)
		assert.NoError(t, err)
	})

	t.Run("SeedAdmin with password > 72 bytes returns error", func(t *testing.T) {
		memDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		assert.NoError(t, err)
		_ = AutoMigrate(memDB)

		longPassCfg := &config.Config{
			AdminUsername: "admin2",
			AdminPassword: strings.Repeat("a", 100),
		}

		err = SeedAdmin(memDB, longPassCfg)
		assert.Error(t, err)
	})

	t.Run("SeedAdmin and AutoMigrate on closed db error", func(t *testing.T) {
		memDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		assert.NoError(t, err)

		sqlDB, err := memDB.DB()
		assert.NoError(t, err)
		_ = sqlDB.Close()

		err = AutoMigrate(memDB)
		assert.Error(t, err)

		err = SeedAdmin(memDB, cfg)
		assert.Error(t, err)
	})

	t.Run("InitDB with postgres branch returns error when no server", func(t *testing.T) {
		pgCfg := &config.Config{
			DBConnection: "pgsql",
			DBHost:       "127.0.0.1",
			DBPort:       "59999",
			DBDatabase:   "nonexistent",
			DBUsername:   "invalid",
			DBPassword:   "invalid",
			DBSSLMode:    "disable",
		}
		db, err := InitDB(pgCfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}
