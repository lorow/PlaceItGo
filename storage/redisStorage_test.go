package storage

import (
	"placeitgo/config"
	"placeitgo/model"
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

	entries := []model.ImageDBEntry{
		{
			Author: "lorow",
			Title:  "testing_entry",
			Link:   "google.com",
			Width:  1920,
			Height: 1080,
		},
	}

	err := redisCache.SaveImageEntries(entries, "fox")
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

func Test_fetching_image_direct_key(t *testing.T) {
	setup()
	defer teardown()
	// set the image first

	entries := []model.ImageDBEntry{
		{
			Author: "lorow",
			Title:  "testing_entry",
			Link:   "google.com",
			Width:  1920,
			Height: 1080,
		},
	}
	redisCache.SaveImageEntries(entries, "fox")

	image, err := redisCache.GetImage(1920, 1080, "fox")

	if err != nil {
		t.Errorf("Error while fetching image: %s", err)
	} else {
		if image.Author != "lorow" {
			t.Errorf("Author doesn't match: %s %s", "lorow", image.Author)
		}
		if image.Title != "testing_entry" {
			t.Errorf("Ttile doesn't match: %s %s", "testing_entry", image.Title)
		}
		if image.Link != "google.com" {
			t.Errorf("Link doesn't match: %s %s", "google.com", image.Link)
		}
		if image.Width != 1920 {
			t.Errorf("Width doesn't match: %d %d", 1920, image.Width)
		}
		if image.Height != 1080 {
			t.Errorf("Height doesn't match: %d %d", 1080, image.Height)
		}
	}
}

func Test_fetching_image_similar_key(t *testing.T) {
	setup()
	defer teardown()

	entries := []model.ImageDBEntry{
		{
			Author: "lorow",
			Title:  "testing_entry",
			Link:   "google.com",
			Width:  1940,
			Height: 1100,
		},
	}
	redisCache.SaveImageEntries(entries, "fox")
	image, err := redisCache.GetImage(1920, 1080, "fox")

	if err != nil {
		t.Errorf("Error while fetching image: %s", err)
	} else {
		if image.Author != "lorow" {
			t.Errorf("Author doesn't match: %s %s", "lorow", image.Author)
		}
		if image.Title != "testing_entry" {
			t.Errorf("Ttile doesn't match: %s %s", "testing_entry", image.Title)
		}
		if image.Link != "google.com" {
			t.Errorf("Link doesn't match: %s %s", "google.com", image.Link)
		}
		if image.Width != 1940 {
			t.Errorf("Width doesn't match: %d %d", 1940, image.Width)
		}
		if image.Height != 1100 {
			t.Errorf("Height doesn't match: %d %d", 1100, image.Height)
		}
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
