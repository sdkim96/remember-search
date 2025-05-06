package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

func NewElasticClient(elasticHost string, elasticAPIKey string) (*elasticsearch.Client, error) {

	hosts := strings.Split(elasticHost, ",")
	cfg := elasticsearch.Config{Addresses: hosts, APIKey: elasticAPIKey}

	c, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Bulk(
	client *elasticsearch.Client,
	index string,
	dto *[]CompanyAnalysisDTO,
) error {

	var buf bytes.Buffer

	for _, a := range *dto {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, a.DocumentID, "\n"))
		data, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Cannot encode article %s: %s", a.DocumentID, err)
		}

		// Append newline to the data payload
		//
		data = append(data, "\n"...) // <-- Comment out to trigger failure for batch
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		res, err := client.Bulk(bytes.NewReader(buf.Bytes()), client.Bulk.WithIndex(index))
		if err != nil {
			return err
		}
		defer res.Body.Close()
		fmt.Printf("Bulk Success: %s\n", a.DocumentID)
	}

	return nil
}
