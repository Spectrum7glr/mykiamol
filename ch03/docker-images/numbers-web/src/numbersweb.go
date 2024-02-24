package numbersweb

import (
	"html/template"
	"log"
	"net/http"

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

func RunCLI() {
	viper.AutomaticEnv() // Automatically read environment variables

	// Set default value
	viper.SetDefault("API_ENDPOINT", "http://localhost:8091/rnd")

	mux := http.NewServeMux()
	nw := http.HandlerFunc(numbersweb)
	mux.Handle("/", nw)
	log.Println("Starting server on :80")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}
