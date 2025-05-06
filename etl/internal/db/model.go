package db

import "fmt"

type OfficeDescriptionModel struct {
	RemeberID int    `json:"remember_id"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Address   string `json:"address"`
	Content   string `json:"content"`
}

func (o *OfficeDescriptionModel) GetDescription() string {

	return fmt.Sprintf(`
	회사 이름: %s
	주소: %s
	설명: %s
	`, o.Title, o.Address, o.Content)
}

type CompanyInfoDTO struct {
	CompanySummary  string   `json:"summary" jsonschema_description:"회사에 대한 요약"`
	CompanyKeywords []string `json:"answer" jsonschema_description:"회사에 대한 키워드들"`
}
