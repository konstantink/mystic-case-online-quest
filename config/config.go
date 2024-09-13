package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	env := os.Getenv("ENV")
	src := os.Getenv("PROJECT_FOLDER")

	if env != "" {
		err := godotenv.Overload(fmt.Sprintf("%s/.env.%s", src, env))
		if err != nil {
			fmt.Printf("Error loading %s/.env.%s file", src, env)
		}
	} else {
		err := godotenv.Overload(".env")
		if err != nil {
			fmt.Print("Error loading .env file")
		}
	}

	return os.Getenv(key)
}
