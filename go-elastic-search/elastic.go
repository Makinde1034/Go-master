package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

var es, _ = elasticsearch.NewDefaultClient()

func Exit() {
	fmt.Println("Goodbye!")
	os.Exit(0)
}

func ReadText(reader *bufio.Scanner, prompt string) string {
	fmt.Print(prompt + ": ")
	reader.Scan()
	return reader.Text()
}

func LoadData() {
	var spaceCrafts []map[string]interface{}
	pageNumber := 0

	for {
		response, err := http.Get("http://stapi.co/api/v1/rest/spacecraft/search?pageSize=100&pageNumber=" + strconv.Itoa(pageNumber))

		if err != nil {
			log.Fatal(err)
		}
		body, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		// page := result["page"].(map[string]interface{})
		// totalPages := int(page["totalPages"].(float64))

		crafts := result["spacecrafts"].([]interface{})

		for _, c := range crafts {
			spaceCrafts = append(spaceCrafts, c.(map[string]interface{}))
		}

		pageNumber++

		if pageNumber >= 3 {
			break
		}

	}

	for _, data := range spaceCrafts {
		uid, _ := data["uid"].(string)
		jsonString, _ := json.Marshal(data)
		request := esapi.IndexRequest{Index: "stsc", DocumentID: uid, Body: strings.NewReader(string(jsonString))}
		request.Do(context.Background(), es)
	}

	fmt.Println(len(spaceCrafts))

}

func search(reader *bufio.Scanner, queryType string) {
	key := ReadText(reader, "Enter key")
	value := ReadText(reader, "Enter value")

	var buffer bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			queryType: map[string]interface{}{
				key: value,
			},
		},
	}

	json.NewEncoder(&buffer).Encode(query)
	response, _ := es.Search(es.Search.WithIndex("stsc"), es.Search.WithBody(&buffer))
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	fmt.Println(result)
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("0) Exit")
		fmt.Println("1) Load spacecraft")
		fmt.Println("2) Get spacecraft")
		fmt.Println("3) Match search")
		option := ReadText(reader, "Enter option")
		if option == "0" {
			Exit()
		} else if option == "1" {
			LoadData()

		} else if option == "3" {
			search(reader, "match")
		} else {
			fmt.Println("Invalid option")
		}
	}
}
