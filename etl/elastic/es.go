package elastic

import (
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
