package initializers

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	PORT       string `mapstructure:"PORT"`
	DB_URL     string `mapstructtur:"DB_URL"`
	JWT_SECRET string `mapstructure:"JWT_SECRET"`
}

func LoadEnv() error {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	requiredKeys := getRequiredKeys(config)
	missingKeys := checkMissingKeys(requiredKeys, config)

	if len(missingKeys) > 0 {
		return fmt.Errorf("following environment variables not found: %v", missingKeys)
	}

	return nil
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
