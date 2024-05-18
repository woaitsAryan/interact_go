package initializers

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Environment string

const (
	DevelopmentENV Environment = "development"
	ProductionENV  Environment = "production"
)

type Config struct {
	PORT                 string      `mapstructure:"PORT"`
	ENV                  Environment `mapstructure:"ENV"`
	DB_HOST              string      `mapstructure:"DB_HOST"`
	DB_PORT              string      `mapstructure:"DB_PORT"`
	DB_NAME              string      `mapstructure:"DB_NAME"`
	DB_USER              string      `mapstructure:"DB_USER"`
	DB_PASSWORD          string      `mapstructure:"DB_PASSWORD"`
	REDIS_HOST           string      `mapstructure:"REDIS_HOST"`
	REDIS_PORT           string      `mapstructure:"REDIS_PORT"`
	REDIS_PASSWORD       string      `mapstructure:"REDIS_PASSWORD"`
	JWT_SECRET           string      `mapstructure:"JWT_SECRET"`
	EARLY_ACCESS_SECRET  string      `mapstructure:"EARLY_ACCESS_SECRET"`
	FRONTEND_URL         string      `mapstructure:"FRONTEND_URL"`
	BACKEND_URL          string      `mapstructure:"BACKEND_URL"`
	ML_URL               string      `mapstructure:"ML_URL"`
	API_TOKEN            string      `mapstructure:"API_TOKEN"`
	SENDGRID_KEY         string      `mapstructure:"SENDGRID_KEY"`
	GMAIL_KEY            string      `mapstructure:"GMAIL_KEY"`
	GOOGLE_CLIENT_ID     string      `mapstructure:"GOOGLE_CLIENT_ID"`
	GOOGLE_CLIENT_SECRET string      `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GOOGLE_OAUTH_STATE   string      `mapstructure:"GOOGLE_OAUTH_STATE"`
	GCP_PROJECT          string      `mapstructure:"GCP_PROJECT"`
	GCP_PUBLIC_BUCKET    string      `mapstructure:"GCP_PUBLIC_BUCKET"`
	GCP_PRIVATE_BUCKET   string      `mapstructure:"GCP_PRIVATE_BUCKET"`
	GCP_CREDS            string      `mapstructure:"GCP_CREDS"`
	LOGGER_URL           string      `mapstructure:"LOGGER_URL"`
	LOGGER_SECRET        string      `mapstructure:"LOGGER_SECRET"`
	LOGGER_TOKEN         string      `mapstructure:"LOGGER_TOKEN"`
	MAILER_URL           string      `mapstructure:"MAILER_URL"`
	MAILER_SECRET        string      `mapstructure:"MAILER_SECRET"`
	MAILER_TOKEN         string      `mapstructure:"MAILER_TOKEN"`
	POPULATE_DUMMIES     bool        `mapstructure:"POPULATE_DUMMIES"`
	POPULATE_ORGS        bool        `mapstructure:"POPULATE_ORGS"`
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

	if CONFIG.ENV != DevelopmentENV && CONFIG.ENV != ProductionENV {
		err := fmt.Errorf("invalid ENV value: %s", CONFIG.ENV)
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
