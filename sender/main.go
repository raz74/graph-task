package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const receiverUrl = "http://localhost:8000/message"

func main() {
	var wg sync.WaitGroup
	n := 10000
	wg.Add(n)
	for i := 0; i < n; i++ {

		length := rand.Intn(8000-50) + 50
		s := RandStringBytes(length)
		go func() {
			sendRequest(s)
			wg.Done()
		}()
	}
	wg.Wait()
}

func sendRequest(s string) {

	request, err := http.NewRequest("POST", receiverUrl, bytes.NewBuffer([]byte(s)))
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

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
