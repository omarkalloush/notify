package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

type Repository struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannels := os.Getenv("SLACK_CHANNELS")
	repositoriesJSON := os.Getenv("repositoriesJSON")

	var repositories []Repository
	err = json.Unmarshal([]byte(repositoriesJSON), &repositories)
	if err != nil {
		fmt.Println("Error parsing repositories JSON:", err)
		return
	}

	api := slack.New(slackToken)

	for {
		for _, repo := range repositories {
			latestRelease, err := getLatestRelease(repo.Owner, repo.Name)
			if err != nil {
				fmt.Println("Error getting latest release:", err)
				continue
			}

			message := fmt.Sprintf("New release for %s/%s: %s (%s)", repo.Owner, repo.Name, latestRelease.Name, latestRelease.PublishedAt)
			fmt.Println(message)

			channelID, timestamp, err := api.PostMessage(slackChannels, slack.MsgOptionText(message, false))
			if err != nil {
				fmt.Println("Error posting message:", err)
			} else {
				fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
			}
		}

		time.Sleep(60 * time.Minute)
	}
}

type GitHubRelease struct {
	Name        string `json:"name"`
	PublishedAt string `json:"published_at"`
}

func getLatestRelease(owner, repo string) (*GitHubRelease, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get latest release, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release GitHubRelease
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}
