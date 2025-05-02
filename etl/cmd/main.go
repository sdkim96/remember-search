package main

import (
	"fmt"

	"github.com/sdkim96/remember-search/etl/internal"
)

func main() {

	fmt.Printf("Load Settings...\n")

	settings := internal.GetSettings()
	dbHandler := internal.InitDB(settings.GetPGURL())
	defer dbHandler.Close()

	fmt.Printf("Check DB Health...\n")
	err := dbHandler.GetDBHealth()
	if err != nil {
		fmt.Printf("DB Health Check Failed: %v\n", err)
	} else {
		fmt.Printf("DB Health Check Passed\n")
	}

	dbHandler.GetUsers()

}
