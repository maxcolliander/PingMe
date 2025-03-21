package utils

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
)

type ThreadsClient struct {
	ID string `json:"id"`
}

func ConnectThreads(accessToken string) (*ThreadsClient, error) {
	url := "https://graph.threads.net/v1.0/me?fields=id,name"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating the request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making the request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading the response: %w", err)
	}

	var client ThreadsClient
	if err := json.Unmarshal(body, &client); err != nil {
		return nil, fmt.Errorf("Error parsing JSON response: %w", err)
	}

	fmt.Printf("Connected to Threads: ID = %v\n", client.ID)
	return &client, nil
}

func PostToThreads(client *ThreadsClient, accessToken, message string) (bool, error) {
	url := fmt.Sprintf("https://graph.threads.net/v1.0/%s/threads", client.ID)
	postData := map[string]interface{}{
		"text": 	message,
		"media_type": "text",
	}

	jsonData, err := json.Marshal(postData)
	if err != nil {
		return false, fmt.Errorf("Error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("Error creating post request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("Error making the request to Threads: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("Error reading the response: %w", err)
	}

	log.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Failed to publish media container. Status code: %d, Response: %s", resp.StatusCode, string(body))
	}

	var mediaResponse map[string]interface{}
	if err := json.Unmarshal(body, &mediaResponse); err != nil {
		return false, fmt.Errorf("Error unmarshaling response: %w", err)
	}

	creationID, ok := mediaResponse["id"].(string)

	if !ok {
		return false, fmt.Errorf("Missing media ID in response")
	}

	publishSuccess, publishErr := PublishThreadsMedia(client, accessToken, creationID)
    if publishErr != nil {
        return false, fmt.Errorf("Error publishing thread: %w", publishErr)
    }

    if publishSuccess {
        log.Println("Successfully published the thread.")
        return true, nil
    }

    return false, nil
}

func PublishThreadsMedia(client *ThreadsClient, accessToken, creationID string) (bool, error) {
	url := fmt.Sprintf("https://graph.threads.net/v1.0/%s/threads_publish", client.ID)
	params := fmt.Sprintf("creation_id=%s&access_token=%s", creationID, accessToken)
    fullURL := fmt.Sprintf("%s?%s", url, params)
	req, err := http.NewRequest("POST", fullURL, nil)
    if err != nil {
        return false, fmt.Errorf("Error creating the request to publish media container: %w", err)
    }
	resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return false, fmt.Errorf("Error publishing the media container: %w", err)
    }
    defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return false, fmt.Errorf("Error reading the response: %w", err)
    }

	log.Printf("Response Body from Publishing: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
        return false, fmt.Errorf("Failed to publish media container. Status code: %d, Response: %s", resp.StatusCode, string(body))
    }

	var publishResponse map[string]interface{}
    if err := json.Unmarshal(body, &publishResponse); err != nil {
        return false, fmt.Errorf("Error unmarshaling publish response: %w", err)
    }

	creationID, ok := publishResponse["id"].(string)
    if !ok {
        return false, fmt.Errorf("Missing media ID in response")
    }

	log.Printf("Successfully published media container. Media ID: %s\n", creationID)
    return true, nil
}