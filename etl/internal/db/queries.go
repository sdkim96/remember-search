package db

import (
	"context"
	"log"
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
INSERT INTO elastic_index_meta
  ( remember_id
  , index_name
  , document_id
  , status
  , error_message
  , created_at
  , updated_at
  , load_cnt
  )
VALUES
  ( $1
  , $2
  , $3
  , 'success'
  , NULL
  , now()
  , now()
  , $4
  );
INSERT INTO elastic_index_content (
	document_id, summary, tags
) VALUES (
	$3, $5, $6
)
;
`

const getMaxCntSQL string = `
SELECT MAX(load_cnt) AS max_load_cnt
FROM elastic_index_meta
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

func (h *DBHandler) InsertESContent(dtos *[]elastic.CompanyAnalysisDTO, indexName string) error {

	var (
		maxLoadCnt int
		ctx        context.Context
	)

	ctx = context.Background()

	row := h.conn.QueryRowContext(
		ctx,
		getMaxCntSQL,
	)
	err := row.Scan(&maxLoadCnt)
	if err != nil {
		log.Printf("Error getting max load count: %v", err)
		maxLoadCnt = 999999999
	}

	maxLoadCnt++

	stmt, err := h.conn.Prepare(insertESContentSQL)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, dto := range *dtos {
		if _, err := stmt.Exec(
			dto.RemeberID,
			indexName,
			dto.DocumentID,
			maxLoadCnt,
			dto.Summary,
			dto.Tags,
		); err != nil {
			log.Printf("Error executing statement: %v", err)
		}
	}

	return nil
}
