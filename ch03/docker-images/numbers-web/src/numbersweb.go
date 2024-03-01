package numbersweb

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

func numbersweb(w http.ResponseWriter, r *http.Request) {
	apiEndpoint := viper.GetString("API_ENDPOINT")
	t, err := template.New("webpage").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, map[string]string{
		"ApiEndpoint": apiEndpoint,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type APIResponse struct {
	Data string `json:"data"`
}

type Client struct {
	APIEndpoint *url.URL
	HTTPClient  *http.Client
}

func NewClient(apiEndpoint string) *Client {
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	URL := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}

	return &Client{
		APIEndpoint: URL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func makeRequest(c *Client) (string, error) {

	req, err := http.NewRequest("GET", c.APIEndpoint.String(), nil)
	log.Println(c.APIEndpoint.String())
	if err != nil {
		return "", err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var jsonresp APIResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonresp)
	if err != nil {
		return "", err
	}
	return jsonresp.Data, nil
}

func numbersweb2(w http.ResponseWriter, r *http.Request) {
	apiEndpoint := viper.GetString("API_ENDPOINT")
	c := NewClient(apiEndpoint)
	data, err := makeRequest(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(data))
}

func RunCLI() {
	viper.AutomaticEnv() // Automatically read environment variables

	// Set default value
	viper.SetDefault("API_ENDPOINT", "http://localhost:8091/rnd")

	mux := http.NewServeMux()
	nw := http.HandlerFunc(numbersweb)
	mux.Handle("/rnd", nw)
	nw2 := http.HandlerFunc(numbersweb2)
	mux.Handle("/rnd2", nw2)
	log.Println("Starting server on :80")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}
