package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const receiverUrl = "http://localhost:8000/message"

func main() {
	for i := 0; i < 10000; i++ {
		length := rand.Intn(8000-50) + 50
		s := RandStringBytes(length)
		go sendRequest(s)
	}

	time.Sleep(time.Second * 10)
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
