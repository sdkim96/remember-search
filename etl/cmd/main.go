package main

import (
	"github.com/sdkim96/remember-search/etl/internal"
	"github.com/sdkim96/remember-search/etl/internal/db"
	"github.com/sdkim96/remember-search/etl/pipeline"
)

func main() {

	// 1. Initialize application settings
	settings := internal.GetSettings()
	defaultAuthor := settings.GetAuthor()
	openAIAPIMaxQuotas := settings.GetOpenAIAPIMaxQuotas()
	elasticHost := settings.GetElasticHost()
	elasticAPIKey := settings.GetElasticAPIKey()

	// 2. Initialize database driver
	dbHandler := db.InitDB(settings.GetPGURL())
	defer dbHandler.Close()

	// 3. Check and ping the database
	dbHandler.GetDBHealth()

	// 4. Run the ETL process
	//    - Early Part: Extract and Transform
	earlyPipeline := &pipeline.EarlyPart{
		Invoker: defaultAuthor,
	}
	err := pipeline.Execute(earlyPipeline, dbHandler)
	if err != nil {
		panic(err)
	}
	//    - Late Part: Load
	latePipeline := &pipeline.LatePart{
		Invoker:            defaultAuthor,
		OpenAIAPIMaxQuotas: openAIAPIMaxQuotas,
		ElasticHost:        elasticHost,
		ElasticAPIKey:      elasticAPIKey,
	}
	lateError := pipeline.Execute(latePipeline, dbHandler)
	if lateError != nil {
		panic(lateError)
	}
}
