package leaderboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	events "github.com/prestonchoate/space-shmup/Systems/Events"
	events_data "github.com/prestonchoate/space-shmup/Systems/Events/Data"
)

type Client struct {
	endpoint    string
	apiKey      string
	authToken   string
	tokenExpiry time.Time
}

type clientLoginRequest struct {
	ApiKey string `json:"api_key"`
}

type clientLoginResponse struct {
	Token     string    `json:"token"`
	Exiration time.Time `json:"expires_at"`
}

type submissionRequest struct {
	Version  string `json:"version"`
	Initials string `json:"initials"`
	Score    int64  `json:"score"`
}

type submissionResponse struct {
	Id        int64     `json:"id"`
	Version   string    `json:"version"`
	Initials  string    `json:"initials"`
	Score     int64     `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	Version string
	apiKey  string
)

func NewLeaderboardClient() *Client {
	log.Println("Leaderboard Client: Creating new client with version: ", Version)
	if Version == "" {
		Version = "debug"
	}

	if apiKey == "" {
		if val, ok := os.LookupEnv("API_KEY"); ok {
			apiKey = val
		} else {
			log.Println("Leaderboard Client: no api key present. Unable to send high scores")
		}
	}
	c := Client{
		endpoint: "https://leaderboards.prestonchoate.dev/api",
		apiKey:   apiKey,
	}

	return &c
}

func (c *Client) HandleHighScoreSubmission(e events.Event) {
	if data, ok := e.Data.(events_data.HighScoreData); ok {
		log.Println("Leaderboard Client: Submitting high score")
		if c.tokenExpiry.Before(time.Now()) {
			err := c.refreshAuthToken()
			if err != nil {
				log.Println("Leaderboard Client: Could not submit score")
				events.GetEventManagerInstance().Emit(events_data.AddMessage, events_data.AddMessageData{
					Message: "Failed to upload score " + err.Error(),
					Timer:   2.0,
					Color:   rl.Red,
				})
				events.GetEventManagerInstance().Emit(events_data.ScoreSubmissionComplete, events_data.ScoreSubmissionCompleteData{
					Success: false,
				})
				return
			}
		}
		err := c.submitScore(submissionRequest{
			Version:  Version,
			Initials: data.Initials,
			Score:    data.Score,
		})

		if err != nil {
			log.Println("Leaderboard Client: Could not submit score")
			events.GetEventManagerInstance().Emit(events_data.AddMessage, events_data.AddMessageData{
				Message: "Failed to upload score " + err.Error(),
				Timer:   2.0,
				Color:   rl.Red,
			})
			events.GetEventManagerInstance().Emit(events_data.ScoreSubmissionComplete, events_data.ScoreSubmissionCompleteData{
				Success: false,
			})
			return
		}

		events.GetEventManagerInstance().Emit(events_data.AddMessage, events_data.AddMessageData{
			Message: "Score Uploaded!",
			Timer:   2.0,
			Color:   rl.DarkGreen,
		})
		events.GetEventManagerInstance().Emit(events_data.ScoreSubmissionComplete, events_data.ScoreSubmissionCompleteData{
			Success: true,
		})
	}
}

func (c *Client) submitScore(submission submissionRequest) error {
	jsonStr, err := json.Marshal(submission)
	if err != nil {
		log.Println("Leaderboard Client: failed to marshal submission request ", err)
		return fmt.Errorf("failed to create submission payload")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/scores", c.endpoint), bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("Leaderboard Client: failed to create http request ", err)
		return fmt.Errorf("failed to create http client")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.authToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Leaderboard Client: failed to send request ", err)
		return fmt.Errorf("request failure")
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Leaderboard Client: failed to ready request body ", err)
		return fmt.Errorf("failed to parse response")
	}

	if resp.StatusCode != http.StatusCreated {
		log.Println("Leaderboard Client: submission error ", string(data))
		return fmt.Errorf("%v", string(data))
	}

	var scoreResp submissionResponse
	err = json.Unmarshal(data, &scoreResp)
	if err != nil {
		log.Println("Leaderboard Client: failed to unmarshal submission response ", err)
		return fmt.Errorf("failed to parse response data")
	}

	log.Printf("Leaderboard Client: Submission success. New score data: %+v\n", scoreResp)

	return nil
}

func (c *Client) refreshAuthToken() error {
	if c.apiKey == "" {
		return fmt.Errorf("no api key")
	}

	loginReq := clientLoginRequest{
		ApiKey: c.apiKey,
	}

	jsonStr, err := json.Marshal(loginReq)
	if err != nil {
		log.Println("Leaderboard Client: failed to marshal login data ", err)
		return fmt.Errorf("failed to create login data")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/auth/client/login", c.endpoint), bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Println("Leaderboard Client: failed to create http request ", err)
		return fmt.Errorf("failed to create login request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Leaderboard Client: failed to send request ", err)
		return fmt.Errorf("failed to send login request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Leaderboard Client: login error ", resp.Status)
		return fmt.Errorf("login error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Leaderboard Client: failed to read response body ", err)
		return fmt.Errorf("failed to read response body")
	}

	var loginResp clientLoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		log.Println("Leaderboard Client: failed to parse response body ", err)
		return fmt.Errorf("failed to parse login response body")
	}

	c.authToken = loginResp.Token
	c.tokenExpiry = loginResp.Exiration
	return nil
}
