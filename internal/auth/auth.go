package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/sverdejot/espotifai/internal/model"
)

type SpotifyAuth struct {
	clientId	string
	clientSecret	string
}

func NewSpotifyAuth(clientId, clientSecret string) *SpotifyAuth {
	return &SpotifyAuth{
		clientId: clientId,
		clientSecret: clientSecret,
	}
}

func (s SpotifyAuth) AuthEndpoint() string {
	endpoint := url.URL{
		Scheme: "https",
		Host:   "accounts.spotify.com",
		Path:   "authorize",
	}

	query := url.Values{
		"response_type": {"code"},
		"client_id":     {s.clientId},
		"redirect_uri":  {"http://localhost:8080/callback"},
	}

	endpoint.RawQuery = query.Encode()

	return endpoint.String()
}

func (s SpotifyAuth) InitAuth() (string, error) {
	ch := make(chan string)
	spinupServer(ch)
	openBrowser(s.AuthEndpoint())
	return <-ch, nil
}

func spinupServer(ch chan<- string) {
	fmt.Println("spinning up server")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /callback", func(w http.ResponseWriter, r *http.Request) {
		codeParam, ok := r.URL.Query()["code"]
		if !ok || len(codeParam) == 0 {
			panic("cannot get params")
		}
		ch <- codeParam[0]
		json.NewEncoder(w).Encode("<h1>You can now close this tab</h1>")
	})
	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil || err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	waitForServer()
}

func (s SpotifyAuth) RequestToken(code string) (string, error) {
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
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", b64encode(fmt.Sprintf("%s:%s", s.clientId, s.clientSecret))))

	res, err := http.DefaultClient.Do(req)
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

	return response.AccessToken, nil
}

func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func b64encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func waitForServer() {
	timeout := time.NewTimer(10 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timeout.C:
			panic("cannot check status")
		case <-ticker.C:
			res, _ := http.DefaultClient.Get("http://localhost:8080/status")
			if res.Status == "200 OK" {
				return
			}
		}
	}
}

func (c *SpotifyAuth) Me(token string) model.Profile {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var profile model.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		log.Fatal("error while unmarshalling: ", err)
	}

	return profile
}
