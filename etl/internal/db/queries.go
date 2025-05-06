package db

import (
	"context"
	"time"

	"github.com/sdkim96/remember-search/etl/elastic"
)

// get info of remeber
const getOfficeDescriptionSQL string = `
SELECT 
	r.id,
	o.id, 
	o.name as title, 
	o.address, 
	o.description as content 
FROM office o 
LEFT JOIN remeber r
ON r.id = o.id
WHERE o.description NOTNULL
;
`

const insertESContentSQL string = `
INSERT INTO elastic_index_meta (
	remember_id,
	index_name,
	document_id,
	status,
	error_message,
	created_at,
	updated_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7

)
;
INSERT INTO elastic_index_content (
	document_id, title, summary, tags
) VALUES (
	$8, $9, $10, $11
)
;
`

func (h *DBHandler) GetOffices(limit ...int) ([]*OfficeDescriptionModel, error) {

	offices := make([]*OfficeDescriptionModel, 0)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*10,
	)
	defer cancel()
	rows, err := h.conn.QueryContext(ctx, getOfficeDescriptionSQL)
	if err != nil {
		return offices, err
	}
	defer rows.Close()

	for rows.Next() {
		office := &OfficeDescriptionModel{}
		err := rows.Scan(
			&office.RemeberID,
			&office.ID,
			&office.Title,
			&office.Address,
			&office.Content,
		)
		if err != nil {
			return offices, err
		}
		offices = append(offices, office)
	}

	if limit != nil && len(offices) > limit[0] {
		offices = offices[:limit[0]]
	}

	return offices, nil
}

func (h *DBHandler) InsertESContent(dto *elastic.CompanyAnalysisDTO, indexName string) error {

	h.conn.ExecContext(
		context.Background(),
		insertESContentSQL,
		dto.RemeberID,
		indexName,
		dto.DocumentID,
	)
	return nil
}

// INSERT INTO elastic_index_meta (
// 	remember_id,
// 	index_name,
// 	document_id,
// 	status,
// 	error_message,
// 	created_at,
// 	updated_at
// ) VALUES (
// 	$1, $2, $3, $4, $5, $6, $7

// )
// ;
// INSERT INTO elastic_index_content (
// 	document_id, title, summary, tags
// ) VALUES (
// 	$8, $9, $10, $11
// )
// ;
