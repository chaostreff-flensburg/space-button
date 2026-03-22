package main

import (
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
	Matrix   string `json:"matrix,omitempty"`
	Signal   string `json:"signal,omitempty"`
}

type Location struct {
	Address     string  `json:"address,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Timezone    string  `json:"timezone,omitempty"`
	CountryCode string  `json:"country_code,omitempty"`
}

type Feed struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url"`
}

type Feeds struct {
	Blog Feed `json:"blog,omitempty"`
}

type State struct {
	Open       bool  `json:"open"`
	Lastchange int64 `json:"lastchange"`
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

const stateFile string = "state.json"

func readState() State {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return State{Open: false, Lastchange: 0}
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{Open: false, Lastchange: 0}
	}
	return state
}

func writeState(state State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}
	return os.WriteFile(stateFile, data, 0644)
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
			Address:     os.Getenv("ADDRESS"),
			Lat:         lat,
			Lon:         lon,
			Timezone:    os.Getenv("TIMEZONE"),
			CountryCode: os.Getenv("COUNTRY_CODE"),
		},
		Contact: Contact{
			Email:    os.Getenv("EMAIL"),
			Mastodon: os.Getenv("MASTODON"),
			Matrix:   os.Getenv("MATRIX"),
			Signal:   os.Getenv("SIGNAL"),
		},
		Feeds: Feeds{
			Blog: Feed{
				Type: "rss",
				URL:  os.Getenv("BLOG"),
			},
		},
		State:  readState(),
		ExtCCC: os.Getenv("EXT_CCC"),
	}
}

func handleCloseOrOpen(open bool, writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(request.Body)

	if string(body) != os.Getenv("TOKEN") {
		http.Error(writer, "Invalid token", http.StatusUnauthorized)
		return
	}

	if err := writeState(State{Open: open, Lastchange: time.Now().Unix()}); err != nil {
		http.Error(writer, "Failed to update state", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file. Only using Docker environment variables")
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(renderResponse()); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/open", func(writer http.ResponseWriter, request *http.Request) {
		handleCloseOrOpen(true, writer, request)
	})

	http.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) {
		handleCloseOrOpen(false, writer, request)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
