package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func LoadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	return string(content), err
}

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func apiToken() string {
	homeDir, dirErr := os.UserHomeDir()
	CheckError(dirErr)

	tokenPath := homeDir + "/.config/srht/pst-token"
	token, tokenErr := os.ReadFile(tokenPath)

	if tokenErr != nil {
		log.Fatalf("Must have existing sr.ht access token (obtain from https://meta.sr.ht/oauth) at %s", tokenPath)
	}

	tokenString := strings.TrimSuffix(string(token), "\n")
	return "token " + tokenString
}

func Request(method string, url string, data interface{}) string {
	jsonData, _ := json.Marshal(data)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", apiToken())
	client := &http.Client{Timeout: time.Second * 10}

	log.WithFields(log.Fields{
		"method": method,
		"url":    url,
	}).Debug("Making request")
	log.Debugf("With data: %s", string(jsonData))

	response, err := client.Do(request)
	CheckError(err)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	log.Debugf("Server response: %s", string(body))

	return string(body)
}
