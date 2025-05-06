package elastic

import (
	"time"
)

type CompanyAnalysisDTO struct {
	RemeberID  int       `json:"remember_id"`
	DocumentID string    `json:"document_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Summary    string    `json:"summary"`
	Vector     []float64 `json:"vector"`
	Tags       []string  `json:"tags"`
	Timestamp  time.Time `json:"timestamp"`
}
