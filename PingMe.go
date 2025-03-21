package main

import (
	"fmt"
	"PingMe/utils"
	"log"
)

func main() {
	utils.LoadEnv()

	accessToken := utils.GetEnv("THREADS_ACCESS_TOKEN")

	ThreadsClient, err := threadsContact(accessToken)

	if err != nil {
		log.Fatalf("Failed to connect to Threads API: %w", err)
	}

	fmt.Printf("User ID: %v\n", ThreadsClient.ID)
	response := deepSeekContact()
	fmt.Println("Message to post:", response)
	success, err := utils.PostToThreads(ThreadsClient, accessToken, response)
	if err != nil {
		log.Fatalf("Error posting to Threads: %v", err)
	}
	if success {
		fmt.Println("User ID: %v, Posted: %t", ThreadsClient.ID, response )
	}
}

func threadsContact(accessToken string) (*utils.ThreadsClient, error) {
	return utils.ConnectThreads(accessToken)
}

func deepSeekContact() string {
	client := utils.NewDeepSeekClient()
	response, err := client.Chat("Complete the following statement: In the Pond I'd be pondering:")
	if err != nil {
		log.Fatalf("Error sending message: %w", err)
	}
	return response
}



