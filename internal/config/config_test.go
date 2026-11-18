package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("LoadConfig with defaults and unset envs", func(t *testing.T) {
		keys := []string{
			"APP_ENV", "APP_PORT", "APP_KEY", "ADMIN_USERNAME", "ADMIN_PASSWORD",
			"DB_CONNECTION", "DB_HOST", "DB_PORT", "DB_DATABASE", "DB_USERNAME",
			"DB_PASSWORD", "DB_SSLMODE", "SEARCH_ENGINE_ENDPOINT", "OPENAI_ENDPOINT",
			"OPENAI_KEY", "OPENAI_MODEL",
		}
		for _, k := range keys {
			os.Unsetenv(k)
		}

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "development", cfg.AppEnv)
		assert.Equal(t, "8081", cfg.AppPort)
		assert.Equal(t, "default-secret-key-32-characters-minimum", cfg.AppKey)
		assert.Equal(t, "admin", cfg.AdminUsername)
		assert.Equal(t, "adminpassword123", cfg.AdminPassword)
		assert.Equal(t, "pgsql", cfg.DBConnection)
		assert.Equal(t, "127.0.0.1", cfg.DBHost)
		assert.Equal(t, "5432", cfg.DBPort)
		assert.Equal(t, "pabwe_p9_blog", cfg.DBDatabase)
		assert.Equal(t, "postgres", cfg.DBUsername)
		assert.Equal(t, "postgres", cfg.DBPassword)
		assert.Equal(t, "disable", cfg.DBSSLMode)
		assert.Equal(t, "https://api.duckduckgo.com/", cfg.SearchEngineEndpoint)
		assert.Equal(t, "https://api.openai.com/v1", cfg.OpenAIEndpoint)
		assert.Equal(t, "", cfg.OpenAIKey)
		assert.Equal(t, "gpt-4o-mini", cfg.OpenAIModel)
	})

	t.Run("LoadConfig with custom envs", func(t *testing.T) {
		customEnvs := map[string]string{
			"APP_ENV":                "production",
			"APP_PORT":               "9090",
			"APP_KEY":                "custom-secret-key-32-characters-minimum",
			"ADMIN_USERNAME":         "superadmin",
			"ADMIN_PASSWORD":         "superpassword123",
			"DB_CONNECTION":          "sqlite",
			"DB_HOST":                "localhost",
			"DB_PORT":                "5433",
			"DB_DATABASE":            ":memory:",
			"DB_USERNAME":            "appuser",
			"DB_PASSWORD":            "secret",
			"DB_SSLMODE":             "require",
			"SEARCH_ENGINE_ENDPOINT": "https://custom-search.com",
			"OPENAI_ENDPOINT":        "https://custom-ai.com/v1",
			"OPENAI_KEY":             "sk-test-key",
			"OPENAI_MODEL":           "gpt-4-turbo",
		}

		for k, v := range customEnvs {
			os.Setenv(k, v)
		}
		defer func() {
			for k := range customEnvs {
				os.Unsetenv(k)
			}
		}()

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "production", cfg.AppEnv)
		assert.Equal(t, "9090", cfg.AppPort)
		assert.Equal(t, "custom-secret-key-32-characters-minimum", cfg.AppKey)
		assert.Equal(t, "superadmin", cfg.AdminUsername)
		assert.Equal(t, "superpassword123", cfg.AdminPassword)
		assert.Equal(t, "sqlite", cfg.DBConnection)
		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, "5433", cfg.DBPort)
		assert.Equal(t, ":memory:", cfg.DBDatabase)
		assert.Equal(t, "appuser", cfg.DBUsername)
		assert.Equal(t, "secret", cfg.DBPassword)
		assert.Equal(t, "require", cfg.DBSSLMode)
		assert.Equal(t, "https://custom-search.com", cfg.SearchEngineEndpoint)
		assert.Equal(t, "https://custom-ai.com/v1", cfg.OpenAIEndpoint)
		assert.Equal(t, "sk-test-key", cfg.OpenAIKey)
		assert.Equal(t, "gpt-4-turbo", cfg.OpenAIModel)
	})

	t.Run("getEnv with empty string", func(t *testing.T) {
		os.Setenv("TEST_EMPTY_VAR", "")
		defer os.Unsetenv("TEST_EMPTY_VAR")

		val := getEnv("TEST_EMPTY_VAR", "fallback_value")
		assert.Equal(t, "fallback_value", val)
	})
}
