package main

import (
	"fmt"
	"io"
	"net/http"
)

type WebSites struct {
	google string
	bing   string
	apple  string
	naver  string
}

type SearchWebResponse struct {
	body   string
	status int
}

var channels chan *SearchWebResponse = make(chan *SearchWebResponse)

func searchWeb(url string, client ...*http.Client) {
	read := "Error reading body"

	var httpClient *http.Client
	if len(client) > 0 && client[0] != nil {
		httpClient = client[0]
	} else {
		httpClient = &http.Client{}
	}

	httpResp, err := httpClient.Get(url)
	if err != nil {
		fmt.Println("Error:", err)

	}
	defer httpResp.Body.Close()
	readBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	}
	read = string(readBytes)
	channels <- &SearchWebResponse{body: read, status: httpResp.StatusCode}
}

func main() {

	var sites *WebSites = &WebSites{
		google: "https://www.google.com",
		bing:   "https://www.bing.com",
		apple:  "https://www.apple.com",
		naver:  "https://www.naver.com",
	}

	for i := 0; i < 4; i++ {
		switch i {
		case 0:
			go searchWeb(sites.google)
		case 1:
			go searchWeb(sites.bing)
		case 2:
			go searchWeb(sites.apple)
		case 3:
			go searchWeb(sites.naver)
		default:
			fmt.Println("Invalid site")
		}
	}

	for i := 0; i < 4; i++ {
		result := <-channels
		fmt.Printf("Response from site %d, Status: %d\n", i, result.status)
	}
	fmt.Println("All requests completed.")

}
