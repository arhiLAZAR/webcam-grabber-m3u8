package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func downloadParts(webcamPrefix, webcamPlaylist, dirname string) {
	var prevPart string

	webcamURL := fmt.Sprintf("%s/%s", webcamPrefix, webcamPlaylist)

	for {
		lastPart := getLastPart(webcamURL)
		if lastPart != prevPart {
			downloadFile(webcamPrefix, lastPart, dirname)
		}

		time.Sleep(time.Second * timeout)
		prevPart = lastPart
	}

}

func getLastPart(webcamURL string) string {
	resp, err := http.Get(webcamURL)
	checkErr(err)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	resp.Body.Close()

	bodyList := strings.Split(string(bodyBytes), "\n")

	var parts []string
	for _, line := range bodyList {
		if !strings.HasPrefix(line, "#") && line != "" {
			parts = append(parts, line)
		}
	}

	lastPart := parts[len(parts)-1]

	return lastPart
}

func downloadFile(webcamPrefix, lastPart, dirname string) {

	url := fmt.Sprintf("%s/%s", webcamPrefix, lastPart)

	unixDate := time.Now().Unix()
	fileName := fmt.Sprintf("%s%s%d_%s", dirname, normalVideoPrefix, unixDate, lastPart)

	resp, err := http.Get(url)
	checkErr(err)

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fileName)
	checkErr(err)
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	checkErr(err)

}
