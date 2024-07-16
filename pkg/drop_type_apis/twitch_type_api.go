package drop_type_apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TwitchTypeApi struct {
	ClientID string
	Token    string
}

func (t *TwitchTypeApi) Search(search string) []ApiSearchResponse {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/search/channels", nil)
	if err != nil {
		return nil
	}

	q := req.URL.Query()
	q.Add("query", search)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Client-ID", t.ClientID)
	req.Header.Set("Authorization", "Bearer "+t.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil
	}

	if len(searchResponse.Data) == 0 {
		log.Printf("No results found for search query %s\n", search)
		return nil
	}

	var results []ApiSearchResponse
	for _, item := range searchResponse.Data {
		results = append(results, ApiSearchResponse{
			Search:      search,
			PicturePath: strings.Replace(item.ThumbnailUrl, "{width}", "320", -1),
			Title:       item.DisplayName,
			Subtitle:    item.Title,
			Content:     "https://twitch.tv/" + item.DisplayName,
		})
	}

	return results
}

func (t *TwitchTypeApi) Init() {
	fmt.Println("Initializing Twitch API")
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	if clientID == "" || clientSecret == "" {
		log.Printf("Twitch API client ID or secret not found in environment variable TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET\n")
		return
	}
	resp, err := http.Post("https://id.twitch.tv/oauth2/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error getting token from Twitch API: %s\n", err)
		return
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Printf("Error decoding token from Twitch API: %s\n", err)
		return
	}

	t.ClientID = clientID
	t.Token = tokenResponse.AccessToken
	fmt.Println("Twitch API initialized")
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Channel struct {
	Id           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type SearchResponse struct {
	Data []Channel `json:"data"`
}
