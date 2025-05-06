package pipeline

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

	offices, err := h.GetOffices(2)
	if err != nil {
		return fmt.Errorf("failed to get offices: %w", err)
	}
	fmt.Printf("Got %d offices\n", len(offices))

	wg := &sync.WaitGroup{}
	companyAnalysisDTOs := make([]elastic.CompanyAnalysisDTO, 0)
	dtoChan := make(chan *elastic.CompanyAnalysisDTO, len(offices))
	sem := make(chan struct{}, p.OpenAIAPIMaxQuotas)

	// 1. Execute goroutines over the offices (OpenAI API)
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

			// Call the OpenAI API to get the summary and keywords
			resp, err := ai.InvokeOpenAI[db.CompanyInfoDTO](systemPrompt, userPrompt, openaiClient)
			if err != nil {
				fmt.Printf("Error invoking LLM: %v\n", err)
				return
			}

			// Get the embedding vector for the summary
			vector, err := ai.GetEmbedding(resp.CompanySummary, openaiClient)
			if err != nil {
				fmt.Printf("Error getting embedding: %v\n", err)
				return
			}
			summaryIterator := strings.Split(resp.CompanySummary, "")
			fmt.Printf("LLM Response: %s...\n", summaryIterator[:20])

			// Send the pointer to the channel
			dtoChan <- &elastic.CompanyAnalysisDTO{
				RemeberID:  o.RemeberID,
				DocumentID: uuid.NewString(),
				Title:      o.Title,
				Content:    o.Content,
				Summary:    resp.CompanySummary,
				Vector:     vector,
				Tags:       resp.CompanyKeywords,
				Timestamp:  time.Now(),
			}

		}(office)

	}

	// Close the goroutine channel after all goroutines are done (asynchronously)
	go func() {
		wg.Wait()
		close(dtoChan)
	}()

	// Recieve the results from the channel
	for dto := range dtoChan {
		companyAnalysisDTOs = append(companyAnalysisDTOs, *dto)
	}

	// 2. Insert the results into DB and ES
	fmt.Printf("Inserting into Es... %d\n", len(companyAnalysisDTOs))
	if err := elastic.Bulk(es, "dev_company_analysis", &companyAnalysisDTOs); err != nil {
		return fmt.Errorf("failed to bulk insert into elasticsearch: %w", err)
	}
	return nil

	// 3. Insert the results into DB

}
