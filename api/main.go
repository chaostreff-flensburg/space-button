package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Contact struct {
	Email    string `json:"email,omitempty"`
	Mastodon string `json:"mastodon,omitempty"`
}

type Location struct {
	Address string  `json:"address,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Lon     float64 `json:"lon,omitempty"`
}

type Feed struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url"`
}

type Feeds struct {
	Blog Feed `json:"blog,omitempty"`
}

type Icon struct {
	Open   string `json:"open,omitempty"`
	Closed string `json:"closed,omitempty"`
}

type State struct {
	Open       bool  `json:"open,omitempty"`
	Lastchange int64 `json:"lastchange,omitempty"`
	Icon       Icon  `json:"icon,omitempty"`
}

type Space struct {
	ApiCompatibility []string `json:"api_compatibility"`
	Space            string   `json:"space"`
	Logo             string   `json:"logo"`
	Url              string   `json:"url"`
	Location         Location `json:"location,omitempty"`
	Contact          Contact  `json:"contact"`
	Feeds            Feeds    `json:"feeds,omitempty"`
	State            State    `json:"state,omitempty"`
	ExtCCC           string   `json:"ext_ccc,omitempty"`
}

func getLastChange() int64 {
	file, err := os.Open("lastchange.txt")
	if err != nil {
		return time.Now().Unix()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		ts, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return time.Now().Unix()
		}
		return ts
	}
	return time.Now().Unix()
}

// updateLastChange writes the current timestamp to file
func updateLastChange() {
	file, err := os.Create("lastchange.txt")
	if err != nil {
		fmt.Println("Error updating lastchange file:", err)
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, time.Now().Unix())
	w.Flush()
}

func renderResponse() Space {
	lat, _ := strconv.ParseFloat(os.Getenv("LAT"), 64)
	lon, _ := strconv.ParseFloat(os.Getenv("LON"), 64)

	return Space{
		ApiCompatibility: []string{"14", "15"},
		Space:            os.Getenv("SPACE"),
		Logo:             os.Getenv("LOGO"),
		Url:              os.Getenv("URL"),
		Location: Location{
			Address: os.Getenv("ADDRESS"),
			Lat:     lat,
			Lon:     lon,
		},
		Contact: Contact{
			Email:    os.Getenv("EMAIL"),
			Mastodon: os.Getenv("MASTODON"),
		},
		Feeds: Feeds{
			Blog: Feed{
				Type: "rss",
				URL:  os.Getenv("BLOG"),
			},
		},
		State: State{
			Open:       getState(),
			Lastchange: getLastChange(),
		},
		ExtCCC: os.Getenv("EXT_CCC"),
	}
}

func getState() bool {
	file, err := os.Open("state.txt")
	if err != nil {
		fmt.Println("Error opening state file:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text() == "true"
	}
	return false
}

func setState(state bool) {
	os.Remove("state.txt")
	file, err := os.Create("state.txt")
	if err != nil {
		fmt.Println("Error creating state file:", err)
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, state)
	w.Flush()

	updateLastChange()
}

func handleCloseOrOpen(open bool, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if string(body) != os.Getenv("TOKEN") {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	setState(open)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(renderResponse())
	})

	http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		handleCloseOrOpen(true, w, r)
	})

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		handleCloseOrOpen(false, w, r)
	})

	http.ListenAndServe(":8080", nil)
}
