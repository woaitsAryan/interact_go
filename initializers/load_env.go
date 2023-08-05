package initializers

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	PORT                    string `mapstructure:"PORT"`
	DB_URL                  string `mapstructure:"DB_URL"`
	JWT_SECRET              string `mapstructure:"JWT_SECRET"`
	ENV                     string `mapstructure:"ENV"`
	BASE_URL                string `mapstructure:"BASE_URL"`
	MAILGUN_PRIVATE_API_KEY string `mapstructure:"MAILGUN_PRIVATE_API_KEY"`
	MAILGUN_PUBLIC_API_KEY  string `mapstructure:"MAILGUN_PUBLIC_API_KEY"`
	MAILGUN_DOMAIN          string `mapstructure:"MAILGUN_DOMAIN"`
}

var CONFIG Config

func LoadEnv() {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&CONFIG)
	if err != nil {
		log.Fatal(err)
	}

	requiredKeys := getRequiredKeys(CONFIG)
	missingKeys := checkMissingKeys(requiredKeys, CONFIG)

	if len(missingKeys) > 0 {
		err := fmt.Errorf("following environment variables not found: %v", missingKeys)
		log.Fatal(err)
	}
}

func getRequiredKeys(config Config) []string {
	requiredKeys := []string{}
	configType := reflect.TypeOf(config)

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag != "" {
			requiredKeys = append(requiredKeys, tag)
		}
	}

	return requiredKeys
}

func checkMissingKeys(requiredKeys []string, config Config) []string {
	missingKeys := []string{}

	configValue := reflect.ValueOf(config)
	for _, key := range requiredKeys {
		value := configValue.FieldByName(key).String()
		if value == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	return missingKeys
}
