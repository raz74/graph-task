package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// HTTP endpoint
	url := "http://localhost:8000/message"

	// JSON body
	jsonBody := []byte(`{
		"title": "Post title",
		"body": "Post description",
		"userId": 1
	}`)

	// Create a HTTP post request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/json")

	//create client
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(string(body))

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

}
