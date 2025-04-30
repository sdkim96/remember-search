package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadDotenv() error {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return err
}

type config struct {
	ProjectName   string
	ProjectAuthor string

	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBDriver   string
}

func GetSettings() (*config, error) {
	err := loadDotenv()
	settings := &config{
		ProjectName:   os.Getenv("PROJECT_NAME"),
		ProjectAuthor: os.Getenv("PROJECT_AUTHOR"),
		DBName:        os.Getenv("DB_NAME"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBDriver:      os.Getenv("DB_DRIVER"),
	}
	return settings, err
}

func (c *config) GetPGURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName
}
