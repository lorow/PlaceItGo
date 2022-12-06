package config

import (
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	// Set up test environment variables
	setup()
	defer teardown()

	// Get configuration
	config, err := GetConfig()
	if err != nil {
		t.Errorf("Unexpected error when getting configuration: %v", err)
	}

	// Check configuration values
	if config.RedisURL != "localhost" {
		t.Errorf("Unexpected RedisURL value: %s", config.RedisURL)
	}
	if config.RedisDatabase != 1 {
		t.Errorf("Unexpected RedisDatabase value: %d", config.RedisDatabase)
	}
	if config.RedisPassword != "redis_password" {
		t.Errorf("Unexpected RedisPassword value: %s", config.RedisPassword)
	}
	if config.RedditUsername != "username" {
		t.Errorf("Unexpected RedditUsername value: %s", config.RedditUsername)
	}
	if config.RedditPassword != "password" {
		t.Errorf("Unexpected RedditPassword value: %s", config.RedditPassword)
	}
	if config.RedditAppID != "client_id" {
		t.Errorf("Unexpected RedditAppID value: %s", config.RedditAppID)
	}
	if config.RedditSecret != "client_secret" {
		t.Errorf("Unexpected RedditSecret value: %s", config.RedditSecret)
	}
}

func setup() {
	os.Setenv("redis_server", "localhost")
	os.Setenv("redis_db", "1")
	os.Setenv("redis_password", "redis_password")
	os.Setenv("reddit_username", "username")
	os.Setenv("reddit_password", "password")
	os.Setenv("reddit_client_id", "client_id")
	os.Setenv("reddit_client_secret", "client_secret")
}

func teardown() {
	os.Unsetenv("redis_server")
	os.Unsetenv("redis_db")
	os.Unsetenv("redis_password")
	os.Unsetenv("reddit_username")
	os.Unsetenv("reddit_password")
	os.Unsetenv("reddit_client_id")
	os.Unsetenv("reddit_client_secret")
}
