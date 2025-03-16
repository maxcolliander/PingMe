package main

import (
	"fmt"
	"PingMe/utils"
	"log"
	"net/http"
	"io/ioutil"
)

func main() {
	utils.LoadEnv()

	accessToken := utils.GetEnv("THREADS_ACCESS_TOKEN")

	url := "https://graph.threads.net/v1.0/me?fields=id,name"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal("Error creating the request", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Error making the request: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading the response body: ", err)
	}

	fmt.Println("Response: ", string(body))
}

