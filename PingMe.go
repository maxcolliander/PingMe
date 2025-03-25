package main

import (
	"fmt"
	"PingMe/utils"
	"log"
	"strings"
	"os"
	"strconv"
	"time"
)

func main() {
	var timeStr string

	if len(os.Args) > 1 {
		timeStr = os.Args[1]
		fmt.Println("Scheduled time provided:", timeStr)
	} else {
		fmt.Println("No scheduled time provided, posting immediately.")
	}

	done := make(chan bool)

	go SchedulePost(timeStr, done)

	<-done

}

func threadsContact(accessToken string) (*utils.ThreadsClient, error) {
	return utils.ConnectThreads(accessToken)
}

func deepSeekContact() string {
	client := utils.NewDeepSeekClient()
	response, err := client.Chat("Complete the following statement: In the Pond I'd be pondering:")
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	// Remove alternative if it exists
	cleanedResponse := removeAlternative(response)
	return cleanedResponse
}

func removeAlternative(response string) string {
	firstIdx := strings.Index(response, `"`)

	if firstIdx == -1 {
		return strings.TrimSpace(response)
	}

	secondIdx := strings.Index(response[firstIdx + 1:], `"`)
	if secondIdx == -1 {
		return strings.TrimSpace(response)
	}
	secondIdx += firstIdx + 1

	response = strings.TrimSpace(response[firstIdx + 1 : secondIdx])

	return strings.ReplaceAll(response, "*", "")
}


func sendPost() {
	utils.LoadEnv()

	accessToken := utils.GetEnv("THREADS_ACCESS_TOKEN")

	ThreadsClient, err := threadsContact(accessToken)

	if err != nil {
		log.Fatalf("Failed to connect to Threads API: %v", err)
	}

	fmt.Printf("User ID: %v\n", ThreadsClient.ID)
	response := deepSeekContact()
	fmt.Println("Message to post:", response)

	postNow(ThreadsClient, accessToken, response)
}

func postNow(client *utils.ThreadsClient, accessToken, response string) {
	fmt.Println("POSTING: ", response)
	// success, err := utils.PostToThreads(client, accessToken, response)
	// if err != nil {
	// 	log.Fatalf("Error posting to Threads: %v", err)
	// }
	// if success {
	// 	fmt.Println("User ID: %v, Posted: %s", client.ID, response )
	// }
}

func SchedulePost(timeStr string, done chan bool) {
	// EXPECTING HHMM, 12:45 --> "1245"
	var scheduledTime time.Time
	now := time.Now()

	if timeStr == "" {
		fmt.Println("No time provided, posting immediately")
		sendPost()
		done <- true
		return
	}
	switch len(timeStr) {
	case 4: // HHMM Format
		hours, err1 := strconv.Atoi(timeStr[:2])
		minutes, err2 := strconv.Atoi(timeStr[2:])

		if err1 != nil || err2 != nil || hours < 0 || hours > 23 || minutes < 0 || minutes > 59 {
			log.Fatalf("Invalid time: %s. Expected HHMM format with valid hours (00-23) and minutes (00-59).", timeStr)
		}

		scheduledTime = time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0, now.Location())
		
		if scheduledTime.Before(now) {
			scheduledTime = scheduledTime.Add(24 * time.Hour)
		}
	case 12: // YYYYMMDDHHMM Format
		year, err1 := strconv.Atoi(timeStr[:4])
		month, err2 := strconv.Atoi(timeStr[4:6])
		day, err3 := strconv.Atoi(timeStr[6:8])
		hours, err4 := strconv.Atoi(timeStr[8:10])
		minutes, err5 := strconv.Atoi(timeStr[10:])

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil ||
			month < 1 || month > 12 || day < 1 || day > 31 || hours < 0 || hours > 23 || minutes < 0 || minutes > 59 {
			log.Fatalf("Invalid date/time: %s. Expected YYYYMMDDHHMM format with valid values.", timeStr)
		}

		scheduledTime = time.Date(year, time.Month(month), day, hours, minutes, 0, 0, now.Location())

		if scheduledTime.Before(now) {
			log.Fatalf("Invalid date/time: %s. The scheduled time must be in the future.", timeStr)
		}

	default:
		log.Fatalf("Invalid time format: %s. Use 'HHMM' (for today/tomorrow) or 'YYYYMMDDHHMM' (for a specific date).", timeStr)
	}

	delay := time.Until(scheduledTime)
	fmt.Printf("Scheduled post at %s (in %v) \n", scheduledTime.Format("2006-01-02 15:04"), delay)

	go func() {
		time.Sleep(delay)
		sendPost()
		fmt.Println("Post successfully sent at", time.Now().Format("2006-01-02 15:04"))
		done <- true
	}()

}

