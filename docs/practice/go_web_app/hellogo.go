package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var appVersion = "1.2" //Default/fallback version
var instanceNum int

func getFrontpage(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Fprintf(w, "Hello, Go! I'm instance %d running version %s at %s\n", instanceNum, appVersion, t.Format("2019-01-02 15:04:05"))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", appVersion)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	instanceNum = rand.Intn(1000)
	http.HandleFunc("/", getFrontpage)
	http.HandleFunc("/health", health)
	http.HandleFunc("/version", getVersion)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
