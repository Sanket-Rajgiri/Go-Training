package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	EnvVariables = make(map[string]any)
	dbTypeMap    = map[string]string{
		"dev":   "mysql",
		"local": "sqlite",
	}
	JWT_SECRET                   string
	JWT_EXPIRY_MINUTES           int
	REFRESH_TOKEN_EXPIRY_MINUTES int
	DB_TYPE                      string
	DB_USER                      string
	DB_PASSWORD                  string
	DB_HOST                      string
	DB_NAME                      string
	COLLECTOR_ENDPOINT           string
	ENV_NAME                     string
)

func LoadEnv() error {
	var ok bool
	ENV_NAME, ok = os.LookupEnv("ENV")
	if !ok {
		ENV_NAME = "local"
		log.Println("setting envName = local")
	}
	DB_TYPE = dbTypeMap[ENV_NAME]
	if ENV_NAME != "local" {
		required := map[string]*string{
			"DB_USER":            &DB_USER,
			"DB_PASSWORD":        &DB_PASSWORD,
			"DB_HOST":            &DB_HOST,
			"DB_NAME":            &DB_NAME,
			"COLLECTOR_ENDPOINT": &COLLECTOR_ENDPOINT,
			"JWT_SECRET":         &JWT_SECRET,
		}
		for key, ref := range required {
			value := os.Getenv(key)
			if value == "" {
				return fmt.Errorf("missing required environment variable: %s", key)
			}
			*ref = value
		}
		optional_vars := map[string]*int{
			"JWT_EXPIRY_MINUTES":           &JWT_EXPIRY_MINUTES,
			"REFRESH_TOKEN_EXPIRY_MINUTES": &REFRESH_TOKEN_EXPIRY_MINUTES,
		}
		for key, ref := range optional_vars {
			val := os.Getenv(key)
			jwtExpiry, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			if jwtExpiry < 1 || jwtExpiry > 60 {
				return fmt.Errorf("invalid %s: must be 1–60", key)
			}
			*ref = jwtExpiry
		}
	}

	return nil
}
