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

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func loadEnvVariables() (config *EnvConfigs) {
	// Set default values
	// viper.SetDefault("api_endpoint", "http://numbers-api/rnd")
	viper.SetDefault("api.endpoint", "http://localhost/rnd")
	viper.SetDefault("port", "80")
	viper.SetDefault("timeout", 2)

	// Bind environment variables
	viper.BindEnv("path")
	viper.BindEnv("stocazzo")

	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	viper.AddConfigPath("./config")

	// Tell viper the name of your file
	viper.SetConfigName("default")

	// Tell viper the type of your file
	viper.SetConfigType("json")
	// viper.SetConfigType("env")

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Print("Error reading env file", err)
	}

	// Viper finds env variables that match any of the existing keys
	// and replaces the values of those keys with the values of the corresponding env variables
	// (keys are teh ones coming from config file or that have a default value set)
	viper.AutomaticEnv()

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		viper.ReadInConfig()
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatal(err)
		}
	})
	return
}

// struct to map env values
type EnvConfigs struct {
	Api struct {
		Endpoint string `mapstructure:"ENDPOINT"`
	} `mapstructure:"API"`
	// ApiEndpoint string `mapstructure:"API_ENDPOINT"`
	Port    string `mapstructure:"PORT"`
	Timeout int    `mapstructure:"TIMEOUT"`
	Path    string `mapstructure:"PATH"`
	Logging string `mapstructure:"LOGGING"`
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

func InitialPage(c *EnvConfigs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func ResponsePage(c *EnvConfigs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		resp, err := internalclient.Get(c.Api.Endpoint, nil)
		if err != nil {
			templ.Execute(w, map[string]string{
				"Answer": err.Error(),
				// "ApiEndpoint": c.ApiEndpoint,
				"ApiEndpoint": c.Api.Endpoint,
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
			// "ApiEndpoint": c.ApiEndpoint,
			"ApiEndpoint": c.Api.Endpoint,
			"Token":       token,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ConfigurationPage(c *EnvConfigs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c.Logging == "Development" {
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
			return
		}
		http.Error(w, "Not allowed", http.StatusForbidden)
	}
}

func RunServer() {

	c := loadEnvVariables()

	mux := http.NewServeMux()
	ip := http.HandlerFunc(InitialPage(c))
	mux.Handle("GET /", ip)
	rp := http.HandlerFunc(ResponsePage(c))
	mux.Handle("POST /", rp)
	cp := http.HandlerFunc(ConfigurationPage(c))
	mux.Handle("GET /config", cp)
	// port := viper.GetString("PORT")
	http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux)
}
