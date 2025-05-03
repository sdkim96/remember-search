package pipeline

import (
	"fmt"

	"github.com/sdkim96/remember-search/etl/internal/db"
)

type Pipeline interface {
	Run(h *db.DBHandler) error
}

func Execute(p Pipeline, h *db.DBHandler) error {
	err := p.Run(h)
	if err != nil {
		return fmt.Errorf("failed to run pipeline: %w", err)
	}
	fmt.Println("Pipeline executed successfully.")
	return nil
}

// 1. ETLPipeLine

type ETLPipeLine struct {
	Invoker string
}

func (p *ETLPipeLine) Run(h *db.DBHandler) error {
	// Implement your logic here
	fmt.Println("Running pipeline...")
	return nil
}
