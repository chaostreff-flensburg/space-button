package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Contact struct {
	Email    string `json:"email"`
	Mastodon string `json:"mastodon"`
}

type Location struct {
	Address string `json:"address"`
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
}

type Feeds struct {
	Blog string `json:"blog"`
}

type Space struct {
	Api              string   `json:"api"`
	ApiCompatibility string   `json:"api_compatibility"`
	Space            string   `json:"space"`
	Logo             string   `json:"logo"`
	Url              string   `json:"url"`
	Location         Location `json:"location"`
	Contact          Contact  `json:"contact"`
	Feeds            Feeds    `json:"feeds"`
	Open             bool     `json:"open"`
	ExtCCC           string   `json:"ext_ccc"`
}

func renderResponse() Space {
	return Space{
		Api:              "0.13",
		ApiCompatibility: "14",
		Space:            os.Getenv("SPACE"),
		Logo:             os.Getenv("LOGO"),
		Url:              os.Getenv("URL"),
		Location: Location{
			Address: os.Getenv("ADDRESS"),
			Lat:     os.Getenv("LAT"),
			Lon:     os.Getenv("LON"),
		},
		Contact: Contact{
			Email:    os.Getenv("EMAIL"),
			Mastodon: os.Getenv("MASTODON"),
		},
		Feeds: Feeds{
			Blog: os.Getenv("BLOG"),
		},
		ExtCCC: os.Getenv("EXT_CCC"),
		Open:   getState(),
	}
}

func getState() bool {
	file, err := os.Open("state.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		if scanner.Text() == "true" {
			fmt.Println(scanner.Text())
			return true
		}
	}
	return false
}

func setState(state bool) {
	os.Remove("state.txt")
	file, err := os.Create("state.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, state)
	w.Flush()
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
		if r.Method == "POST" {
			// Read request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			if string(body) != os.Getenv("TOKEN") {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		setState(true)
	})

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Read request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			if string(body) != os.Getenv("TOKEN") {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		setState(false)
	})

	http.ListenAndServe(":8080", nil)
}
