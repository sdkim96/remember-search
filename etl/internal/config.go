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
	projectName   string
	projectAuthor string

	ctxTimeout int

	dbName     string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbDriver   string
}

func GetSettings() *config {
	err := loadDotenv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	settings := &config{
		projectName:   os.Getenv("PROJECT_NAME"),
		projectAuthor: os.Getenv("PROJECT_AUTHOR"),

		ctxTimeout: 10,

		dbName:     os.Getenv("DB_NAME"),
		dbUser:     os.Getenv("DB_USER"),
		dbPassword: os.Getenv("DB_PASSWORD"),
		dbHost:     os.Getenv("DB_HOST"),
		dbPort:     os.Getenv("DB_PORT"),
		dbDriver:   os.Getenv("DB_DRIVER"),
	}
	return settings
}

func (c *config) GetPGURL() string {
	return "postgres://" + c.dbUser + ":" + c.dbPassword + "@" + c.dbHost + ":" + c.dbPort + "/" + c.dbName
}

func (c *config) GetCtxTimeout() int {
	return c.ctxTimeout
}
