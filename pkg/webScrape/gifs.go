package webScrape

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func RequestGif(searchStr, tenorAPIkey string) string {
	URL := "https://g.tenor.com/v1/search?q=" + searchStr + "&key=" + tenorAPIkey + "&limit=1"

	responseResults := GetResponseResults(URL)

	gifURL := ParseGifResponseForGifURL(responseResults)

	return gifURL
}

func GetResponseResults(url string) map[string]interface{} {
	var responseResults map[string]interface{}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&responseResults)
	if err != nil {
		log.Fatal(err)
	}

	return responseResults
}

func ParseGifResponseForGifURL(responseResults map[string]interface{}) string {
	var gifURL string

	for responseResultsKey, responseResultsValue := range responseResults {
		if responseResultsKey == "results" {
			for key, value := range responseResultsValue.([]interface{})[0].(map[string]interface{}) {
				if key == "media" {
					for k, v := range value.([]interface{})[0].(map[string]interface{}) {
						if k == "gif" {
							for finalKey, finalValue := range v.(map[string]interface{}) {
								if finalKey == "url" {
									gifURL = fmt.Sprintf("%v", finalValue)
								}
							}
						}
					}
				}
			}
		}
	}
	return gifURL
}