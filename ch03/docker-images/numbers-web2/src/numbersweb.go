package numbersweb2

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"text/template"
	"time"

	"github.com/spf13/viper"
)

var EnvConfigs *envConfigs

// We will call this in main.go to load the env variables
func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

func loadEnvVariables() (config *envConfigs) {
	viper.SetDefault("api_endpoint", "http://numbers-api/rnd")
	viper.SetDefault("port", "80")
	viper.SetDefault("timeout", 2)
	viper.BindEnv("path")
	viper.BindEnv("stocazzo")

	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	viper.AddConfigPath("./config")

	// Tell viper the name of your file
	viper.SetConfigName("ch04")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Print("Error reading env file", err)
	}
	viper.AutomaticEnv()
	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}

// struct to map env values
type envConfigs struct {
	ApiEndpoint string `mapstructure:"API_ENDPOINT"`
	Port        string `mapstructure:"PORT"`
	Timeout     int    `mapstructure:"TIMEOUT"`
	Path        string `mapstructure:"PATH"`
}

var validtokens = make(map[string]bool)
var internalclient = InitHttpClient(time.Duration(viper.GetInt("timeout")) * time.Second)

type APIResponse struct {
	Data string `json:"data"`
}

func generateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func initialPage(w http.ResponseWriter, r *http.Request) {
	templ, err := template.New("index").Parse(indexPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := generateToken()
	validtokens[token] = true
	err = templ.Execute(w, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func responsePage(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	// Handle error or redirect as appropriate
	// 	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	token := r.FormValue("token")
	if !validtokens[token] {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}

	delete(validtokens, token)
	templ, err := template.New("result").Parse(resultPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token = generateToken()
	validtokens[token] = true
	// resp, err := internalclient.Get(viper.GetString("API_ENDPOINT"), nil)
	resp, err := internalclient.Get(EnvConfigs.ApiEndpoint, nil)
	if err != nil {
		templ.Execute(w, map[string]string{
			"Answer": err.Error(),
			// "ApiEndpoint": viper.GetString("API_ENDPOINT"),
			"ApiEndpoint": EnvConfigs.ApiEndpoint,
			"Token":       token,
		})

		return
	}
	var apiResponse APIResponse
	err = json.Unmarshal([]byte(resp), &apiResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, map[string]string{
		"Answer": apiResponse.Data,
		// "ApiEndpoint": viper.GetString("API_ENDPOINT"),
		"ApiEndpoint": EnvConfigs.ApiEndpoint,
		"Token":       token,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func configurationPage(w http.ResponseWriter, r *http.Request) {
	templ, err := template.New("config").Parse(configPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	list := []string{}
	for _, key := range viper.AllKeys() {
		list = append(list, key+"="+viper.GetString(key))

	}
	slices.Sort(list)
	err = templ.Execute(w, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RunServer() {
	// viper.AutomaticEnv() // Automatically read environment variables

	// Set default value
	// viper.SetDefault("API_ENDPOINT", "http://localhost:8091/rnd")
	EnvConfigs = loadEnvVariables()
	mux := http.NewServeMux()
	ip := http.HandlerFunc(initialPage)
	mux.Handle("GET /", ip)
	rp := http.HandlerFunc(responsePage)
	mux.Handle("POST /", rp)
	cp := http.HandlerFunc(configurationPage)
	mux.Handle("GET /config", cp)
	// port := viper.GetString("PORT")
	http.ListenAndServe(fmt.Sprintf(":%s", EnvConfigs.Port), mux)
}
