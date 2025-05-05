package pipeline

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openai/openai-go"

	"github.com/sdkim96/remember-search/etl/ai"
	"github.com/sdkim96/remember-search/etl/elastic"
	"github.com/sdkim96/remember-search/etl/internal/db"
)

type ETLPipeline interface {
	Run(h *db.DBHandler) error
}

func Execute(p ETLPipeline, h *db.DBHandler) error {
	err := p.Run(h)
	if err != nil {
		return fmt.Errorf("failed to run pipeline: %w", err)
	}
	fmt.Println("Pipeline executed successfully.")
	return nil
}

// First Step of the ETL process
//
// This Early Partial ETL process is responsible for extracting data from the database,
// transforming it using the OpenAI API, and loading it back into the database.
type EarlyPart struct {
	Invoker string
}

// Second Step of the ETL process
//
// This Late Partial ETL process is responsible for loading the transformed data back into the ElasticSearch.
type LatePart struct {
	Invoker            string
	OpenAIAPIMaxQuotas int
	ElasticHost        string
	ElasticAPIKey      string
}

func (p *EarlyPart) Run(h *db.DBHandler) error {

	return nil
}

func (p *LatePart) Run(h *db.DBHandler) error {
	openaiClient := openai.NewClient()
	es, elErr := elastic.NewElasticClient(p.ElasticHost, p.ElasticAPIKey)
	if elErr != nil {
		return fmt.Errorf("failed to create elasticsearch client: %w", elErr)
	}
	res, err := es.Info()
	if err != nil {
		return fmt.Errorf("failed to get elasticsearch info: %w", err)
	}
	defer res.Body.Close()
	fmt.Printf("Elasticsearch Info: %s\n", res.String())

	offices, err := h.GetOffices(10)
	if err != nil {
		return fmt.Errorf("failed to get offices: %w", err)
	}
	fmt.Printf("Got %d offices\n", len(offices))

	wg := &sync.WaitGroup{}
	companyAnalysisDTOs := make([]elastic.CompanyAnalysisDTO, 0)
	dtoChan := make(chan *elastic.CompanyAnalysisDTO, len(offices))

	// Limit the number of concurrent requests to OpenAI API
	sem := make(chan struct{}, p.OpenAIAPIMaxQuotas)

	// 1. Iterate over the offices
	for _, office := range offices {
		wg.Add(1)
		// Acquire a token from the semaphore
		sem <- struct{}{}

		go func(o *db.OfficeDescriptionModel) {

			// Release the token when the goroutine completes
			// This ensures that the semaphore is released even if an error occurs Because of the `defer`
			defer func() {
				<-sem
				wg.Done()
			}()
			fmt.Printf("Processing office: %s\n", o.Title)

			officeInfo := o.GetDescription()
			systemPrompt := fmt.Sprintf(`
## 역할
당신은 회사 정보를 요약하고 정리하는 전문가입니다.

## 회사 정보
%s
			`, officeInfo)

			userPrompt := "회사에 대해 요약과 키워드를 정리해주세요."

			resp, err := ai.InvokeOpenAI[db.CompanyInfoDTO](systemPrompt, userPrompt, openaiClient)
			if err != nil {
				fmt.Printf("Error invoking LLM: %v\n", err)
				return
			}

			vector, err := ai.GetEmbedding(resp.CompanySummary, openaiClient)
			if err != nil {
				fmt.Printf("Error getting embedding: %v\n", err)
				return
			}
			summaryLogging := strings.Split(resp.CompanySummary, "")
			fmt.Printf("LLM Response: %s...\n", summaryLogging[:20])
			dtoChan <- &elastic.CompanyAnalysisDTO{
				Title:     o.Title,
				Content:   officeInfo,
				Summary:   resp.CompanySummary,
				Vector:    vector,
				Tags:      resp.CompanyKeywords,
				Timestamp: time.Now(),
			}

		}(office)

	}

	go func() {
		wg.Wait()
		close(dtoChan)
	}()

	for dto := range dtoChan {
		companyAnalysisDTOs = append(companyAnalysisDTOs, *dto)
	}
	fmt.Println("Inserting into Es... %d\n", len(companyAnalysisDTOs))
	return nil
}
