package storage

import (
	"placeitgo/config"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

var redisCache *RedisCache
var redisServer *miniredis.Miniredis

func Test_hashing(t *testing.T) {
	setup()
	defer teardown()
	firstLink := "google.com"
	secondLink := "reddit.com"

	firstHash := redisCache.hash(firstLink)
	secondHash := redisCache.hash(secondLink)

	if firstHash == secondHash {
		t.Error("hashes did match while they shouldn't have")
	}
}

func Test_saving_image(t *testing.T) {
	setup()
	defer teardown()

	err := redisCache.SaveImage(1920, 1080, "lorow", "testing_entry", "fox", "google.com")
	if err != nil {
		t.Errorf("saving image data returned error: %s", err)
	}

	result := redisServer.HGet("image:2006368837:1920x1080-fox", "title")
	if result == "" {
		t.Errorf("stored image doesn't match the expected result")
	}

	resolutions, err := redisServer.SMembers("fox_resolutions")
	if err != nil {
		t.Errorf("could not get resultions: %s", err)
	}
	if len(resolutions) == 0 {
		t.Errorf("no resolutions stored")
	}
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func setup() {
	redisServer = mockRedis()
	config := config.Config{
		RedisURL:      redisServer.Addr(),
		RedisDatabase: 0,
		RedisPassword: "",
	}
	var err error
	redisCache, err = NewRedisCache(&config)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	redisCache = nil
}
