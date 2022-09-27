package internal

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	RedisURL       string `env:"redis_url"`
	RedisDatabase  int    `env:"redis_databse"`
	RedisPassword  string `env:"redis_password"`
	RedditUsername string `env:"reddit_username"`
	RedditPassword string `env:"reddit_password"`
	RedditAppID    string `env:"reddit_app_id"`
	RedditSecret   string `env:"reddit_secret"`
}

func GetConfig() (*Config, error) {
	config := Config{}
	reflected_config_type := reflect.TypeOf(config)
	reflected_config_pointer := reflect.ValueOf(&config).Elem()

	for i := 0; i < reflected_config_type.NumField(); i++ {
		struct_field := reflected_config_type.Field(i)
		tag, ok := struct_field.Tag.Lookup("env")
		base_value := ""

		if ok {
			base_value = os.Getenv(tag)
		}

		field := reflected_config_pointer.FieldByName(struct_field.Name)
		if field.CanSet() {
			if field.Kind() == reflect.Int {
				i, err := strconv.Atoi(base_value)
				if err != nil {
					i = 0
				}
				field.SetInt(int64(i))
			}
			if field.Kind() == reflect.String {
				field.SetString(base_value)
			}
		} else {
			fmt.Printf("we can't set anything, weird")
		}
	}

	return &config, nil
}
