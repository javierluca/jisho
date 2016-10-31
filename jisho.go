package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
)

const JISHO_URL string = "http://jisho.org/api/v1/search/words"
const USAGE string = "Usage: jisho $word ($num_results)"
const DEFAULT_NUM_RESULTS int = 5

type JishoRequest interface {
	requestMeaning(word string)
}

type ApiRequest struct {
	client  http.Client
	url     string
	channel chan map[string]interface{}
}

func (r *ApiRequest) requestMeaning(word string) {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		log.Fatal("Error creating the http request", err)
	}
	q := req.URL.Query()
	q.Add("keyword", word)
	req.URL.RawQuery = q.Encode()

	resp, err := (&r.client).Do(req)
	if err != nil {
		log.Fatal("Error with the http request", err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading the http response", err)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(bodyText, &dat); err != nil {
		log.Fatal("Error parsing response", err)
	}
	r.channel <- dat
}

func printResult(result map[string]interface{}, num_results int) {
	data := result["data"].([]interface{})

	var i int = 0
	for _, dat := range data {
		d := dat.(map[string]interface{})
		for _, jp := range d["japanese"].([]interface{}) {
			j := jp.(map[string]interface{})
			color.Yellow("%s = %s", j["word"], j["reading"])
		}

		for _, sense := range d["senses"].([]interface{}) {
			s := sense.(map[string]interface{})
			color.Green("meaning: %s", s["english_definitions"])
		}

		i++
		if i >= num_results {
			break
		}
		color.White("-------------------")
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		color.White(USAGE)
		os.Exit(0)
	}

	word := args[0]

	num_results := DEFAULT_NUM_RESULTS
	if len(args) > 1 {
		i, err := strconv.Atoi(args[1])
		if err != nil {
			color.Red("$num_results is not a valid number")
			color.White(USAGE)
			os.Exit(0)
		}
		num_results = i
	}

	jC := make(chan map[string]interface{})
	api := &ApiRequest{http.Client{}, JISHO_URL, jC}

	go api.requestMeaning(word)
	result := <-api.channel
	printResult(result, num_results)
}
