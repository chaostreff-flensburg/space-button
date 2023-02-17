package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const token string = "secret"

type Contact struct {
	Email    string `json:"email"`
	Mastodon string `json:"mastodon"`
}

type Location struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type Space struct {
	Api              string   `json:"api"`
	ApiCompatibility string   `json:"api_compatibility"`
	Space            string   `json:"space"`
	Logo             string   `json:"logo"`
	Url              string   `json:"url"`
	Location         Location `json:"location"`
	Contact          Contact  `json:"contact"`
	Open             bool     `json:"open"`
}

func renderResponse() Space {
	return Space{
		Api:              "0.13",
		ApiCompatibility: "14",
		Space:            "Chaostreff Flensburg",
		Logo:             "https://chaostreff-flensburg.de/wp-content/uploads/2018/03/ctfl-logo.png",
		Url:              "https://chaostreff-flensburg.de",
		Location: Location{
			Address: "Apenrader Str. 49, 24941 Flensburg",
			Lat:     54.785,
			Lon:     9.437,
		},
		Contact: Contact{
			Email:    "mail@chaostreff-flensburg.de",
			Mastodon: "https://chaos.social/@chaos_fl",
		},
		Open: getState(),
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

      if string(body) != token {
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

      if string(body) != token {
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
