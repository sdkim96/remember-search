package main

import (
	"context"
	"errors"
	"fmt"
	"time"

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

	cause := errors.New("DB connection too slow")

	ctx, cancel := context.WithDeadlineCause(context.Background(), time.Now().Add(1*time.Second), cause)
	defer cancel()
	dbHandler.GetUsersWithCtx(ctx)

}
