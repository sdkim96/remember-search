package main

import (
	"fmt"
	"log"

	"github.com/sdkim96/remember-search/etl/internal"
)

func main() {

	fmt.Printf("Load Settings...\n")

	settings, err := internal.GetSettings()
	if err != nil {
		fmt.Printf("Error loading settings: %v\n", err)
		return
	}

	log.Println("Project Name: ", settings.ProjectName)
	db, err := internal.GetDB(settings.GetPGURL())
	fmt.Printf("DB Connection: %v\n", db)

}
