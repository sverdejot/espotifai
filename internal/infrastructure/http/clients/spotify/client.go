package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/sverdejot/espotifai/internal/model"
)

var tokens map[string]string = make(map[string]string)

type Client struct {
	*http.Client

	clientId     string
	clientSecret string
}

func NewClient(clientId, clientSecret string) *Client {
	// TODO: configure transport to RoundTrip request to automatically add headers
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   3 * time.Second,
	}
	return &Client{
		client,
		clientId,
		clientSecret,
	}
}

func (c *Client) Me(token string) model.Profile {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var profile model.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		log.Fatal("error while unmarshalling: ", err)
	}

	return profile
}

func (c *Client) Artists(token string) model.TopArtists {
	return c.Top(token, "artists")
}

func (c *Client) Tracks(token string) model.TopArtists {
	return c.Top(token, "tracks")
}

func (c *Client) Top(token, typpe string) model.TopArtists {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/"+typpe, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var topArtists model.TopArtists
	if err := json.NewDecoder(resp.Body).Decode(&topArtists); err != nil {
		log.Fatal("error while unmarshalling: ", err)
	}

	return topArtists
}

func (c *Client) RequestToken(code string) (string, error) {
	if token, ok := tokens[code]; ok {
		return token, nil
	}
	endpoint := url.URL{
		Scheme: "https",
		Host:   "accounts.spotify.com",
		Path:   "api/token",
	}

	query := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {"http://localhost:8080/callback"},
	}

	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed creating request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", b64encode(fmt.Sprintf("%s:%s", c.clientId, c.clientSecret))))

	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed requesting code: %w", err)
	}

	var response struct {
		AccessToken  string `json:"access_token"`
		Type         string `json:"token_type"`
		Scope        string `json:"scope"`
		Expires      int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed unmarshalling response: %w", err)
	}

	tokens[code] = response.AccessToken

	return response.AccessToken, nil
}

func b64encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
