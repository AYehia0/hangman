package game

import (
	"encoding/json"
	"net/http"
)

const RAND_API_URL = "https://www.wordgamedb.com/api/v1/words/random"

type WordDescription struct {
	Word      string `json:"word"`
	Category  string `json:"category"`
	Length    int    `json:"numLetters"`
	Syllables int    `json:"numSyllables"`
	Hint      string `json:"hint"`
	Guesses   []string
}

func GetRandomWord() (WordDescription, error) {
	var wordInfo WordDescription

	// Send a GET request to the API URL
	resp, err := http.Get(RAND_API_URL)
	if err != nil {
		return wordInfo, err
	}
	defer resp.Body.Close()

	// Check if the response status code is not OK (200)
	if resp.StatusCode != http.StatusOK {
		return wordInfo, err
	}

	// Parse the JSON response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wordInfo); err != nil {
		return wordInfo, err
	}

	return wordInfo, nil

}
