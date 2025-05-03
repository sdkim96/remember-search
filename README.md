# Remeber-Search
This is a Personal Project that contains ETL pipeline and MCP server.

## Overview
My ETL pipeline contains 2 major steps.
1. Extract source data from perplexity, transform to given form of table, Load to PostgreSQL.
2. Index that data to ElasticSearch.

## ETL
Below is tech stacks used in ETL pipeline.

- GoLang, the best language to implement cocurrency.
- Perplexity API
- OpenAI API
- PostgreSQL
- ElasticSearch

### QuickStart
```sh
cd remember-search/etl
go mod init github.com/sdkim96/remember-search/etl
```