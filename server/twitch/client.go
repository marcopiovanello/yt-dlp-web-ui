package twitch

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const twitchAPIURL = "https://api.twitch.tv/helix"

type Client struct {
	authenticationManager AuthenticationManager
}

func NewTwitchClient(am *AuthenticationManager) *Client {
	return &Client{
		authenticationManager: *am,
	}
}

type streamResp struct {
	Data []struct {
		ID        string `json:"id"`
		UserName  string `json:"user_name"`
		Title     string `json:"title"`
		GameName  string `json:"game_name"`
		StartedAt string `json:"started_at"`
	} `json:"data"`
}

func (c *Client) doRequest(endpoint string, params map[string]string) ([]byte, error) {
	token, err := c.authenticationManager.GetAccessToken()
	if err != nil {
		return nil, err
	}

	reqURL := twitchAPIURL + endpoint
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Client-Id", c.authenticationManager.GetClientId())
	req.Header.Set("Authorization", "Bearer "+token.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) PollStream(channel string, liveChannel chan<- *StreamInfo) error {
	body, err := c.doRequest("/streams", map[string]string{"user_login": channel})
	if err != nil {
		return err
	}

	var sr streamResp
	if err := json.Unmarshal(body, &sr); err != nil {
		return err
	}

	if len(sr.Data) == 0 {
		liveChannel <- &StreamInfo{UserName: channel, IsLive: false}
		return nil
	}

	s := sr.Data[0]
	started, _ := time.Parse(time.RFC3339, s.StartedAt)

	liveChannel <- &StreamInfo{
		ID:        s.ID,
		UserName:  s.UserName,
		Title:     s.Title,
		GameName:  s.GameName,
		StartedAt: started,
		IsLive:    true,
	}

	return nil
}
