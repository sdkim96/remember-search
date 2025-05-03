package pipeline

import (
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	"github.com/sdkim96/remember-search/etl/ai"
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
	Invoker string
}

func (p *EarlyPart) Run(h *db.DBHandler) error {

	return nil
}

func (p *LatePart) Run(h *db.DBHandler) error {
	openaiClient := openai.NewClient()
	wg := &sync.WaitGroup{}

	fmt.Println("Running Early Part of Pipeline, Invoker: %s...", p.Invoker)

	offices, err := h.GetOffices(2)
	if err != nil {
		return fmt.Errorf("failed to get offices: %w\n", err)
	}
	fmt.Printf("Got %d offices\n", len(offices))

	for _, office := range offices {
		wg.Add(1)

		go func(o *db.OfficeDescriptionModel) {
			defer wg.Done()
			fmt.Printf("Processing office: %s\n", o.Title)
			// Invoke LLM
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
			fmt.Printf("LLM Response: %s\n", resp.CompanySummary)
		}(office)

	}
	wg.Wait()
	fmt.Println("All offices processed.")

	return nil
}
