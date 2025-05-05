package internal

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func loadDotenv() error {
	files, readErr := os.ReadDir("../")
	for f := range files {
		fmt.Printf("%s\n", files[f].Name())
	}
	if readErr != nil {
		log.Fatal("Error reading directory")
	}

	loadErr := godotenv.Load("../.env")
	if loadErr != nil {
		log.Fatal("Error loading .env file")
	}
	return loadErr
}

type config struct {
	projectName   string
	projectAuthor string

	dbName     string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string

	openAIAPIKey       string
	openAIAPIMaxQuotas string
}

func GetSettings() *config {
	err := loadDotenv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	settings := &config{
		projectName:   os.Getenv("PROJECT_NAME"),
		projectAuthor: os.Getenv("PROJECT_AUTHOR"),

		dbName:     os.Getenv("DB_NAME"),
		dbUser:     os.Getenv("DB_USER"),
		dbPassword: os.Getenv("DB_PASSWORD"),
		dbHost:     os.Getenv("DB_HOST"),
		dbPort:     os.Getenv("DB_PORT"),

		openAIAPIKey:       os.Getenv("OPENAI_API_KEY"),
		openAIAPIMaxQuotas: os.Getenv("OPENAI_API_MAX_QUOTAS"),
	}
	return settings
}

func (c *config) GetPGURL() string {
	return "postgres://" + c.dbUser + ":" + c.dbPassword + "@" + c.dbHost + ":" + c.dbPort + "/" + c.dbName
}

func (c *config) GetAuthor() string {
	return c.projectAuthor
}

func (c *config) GetOpenAIAPIMaxQuotas() int {
	quotas, err := strconv.Atoi(c.openAIAPIMaxQuotas)
	if err != nil {
		log.Fatalf("Error converting OPENAI_API_MAX_QUOTAS to int: %v", err)
	}
	return quotas
}

func (c *config) GetElasticHost() string {
	return os.Getenv("ELASTIC_HOST")
}
func (c *config) GetElasticAPIKey() string {
	return os.Getenv("ELASTIC_API_KEY")
}
