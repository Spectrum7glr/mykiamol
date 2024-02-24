package whoami

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"golang.org/x/sys/unix"
)

func whoami(w http.ResponseWriter, r *http.Request) {

	uname := unix.Utsname{}
	if err := unix.Uname(&uname); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("new request from:", r.RemoteAddr)
	log.Println("original request:", r.RequestURI)
	fmt.Fprintf(w, "I'm %s, running on %s (%s - %s - %s) with %d CPU(s)",
		string(bytes.Trim(uname.Nodename[:], "\x00")),
		string(bytes.Trim(uname.Sysname[:], "\x00")),
		string(bytes.Trim(uname.Machine[:], "\x00")),
		string(bytes.Trim(uname.Release[:], "\x00")),
		string(bytes.Trim(uname.Version[:], "\x00")),
		runtime.NumCPU())
}
func RunCLI() {
	mux := http.NewServeMux()
	wh := http.HandlerFunc(whoami)
	mux.Handle("/", wh)
	// mux.HandleFunc("/", whoami)
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
