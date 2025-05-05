package elastic

import "time"

type CompanyAnalysisDTO struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Vector    []float64 `json:"vector"`
	Tags      []string  `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
}
