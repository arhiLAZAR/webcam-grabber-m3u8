package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/canhlinh/hlsdl"
)

const (
	webcamURL = "https://1strela.ru/webcams/92e9cb983ec44cac95bc1d0d467a1e48/camera.m3u8"
	// dirname           = "/home/laz/dir_to_trash/video/"
	dirname           = "/ffmpeg/" // check trailing slash!!!
	webcamTimezone    = +3
	speedUpMultiplier = 150

	normalVideoPrefix     = "normal_video_"
	normalVideoSuffix     = ".ts"
	normalVideoCleanAfter = true

	speedUpVideoPrefix     = "fast_video_"
	speedUpVideoSuffix     = ".mp4"
	speedUpVideoCleanAfter = true

	convertAndExit = true
)

func main() {
	var convertDate int64 = time.Now().Unix()

	for {
		normalVideoCount := getLastFile(dirname, normalVideoPrefix, normalVideoSuffix)
		normalVideoCount++
		normalVideoName := fmt.Sprintf("%s%06.f%s", normalVideoPrefix, float64(normalVideoCount), normalVideoSuffix)

		recordLiveStream(webcamURL, dirname, normalVideoName, convertDate)

		if timeToConvert(convertDate, webcamTimezone) {
			convertDate = time.Now().Unix()
			go convertAndUpload(dirname)
		}
	}
}

func timeToConvert(convertDate int64, webcamTimezone int) bool {
	currentDate := time.Now().Unix()
	currentHours, err := strconv.Atoi(time.Now().UTC().Format("15"))
	checkErr(err)

	if currentDate-convertDate > 60*60*1 && currentHours+webcamTimezone == 21 {
		return true
	}

	return false
}

func recordLiveStream(webcamURL, dirname, normalVideoName string, convertDate int64) {
	recorder := hlsdl.NewRecorder(webcamURL, dirname)

	recordedFile, err := recorder.Start(normalVideoName, convertDate, webcamTimezone)
	checkErr(err)

	log.Println("Recorded file at ", recordedFile)
}

func convertAndUpload(dirname string) {
	// Speed up video
	speedUpFilenamePrefMulti := fmt.Sprintf("%sx%d_", speedUpVideoPrefix, speedUpMultiplier)

	speedUpVideoCount := getLastFile(dirname, speedUpFilenamePrefMulti, speedUpVideoSuffix)
	speedUpVideoCount++
	speedUpVideoName := fmt.Sprintf("%s%s%06.f%s", dirname, speedUpFilenamePrefMulti, float64(speedUpVideoCount), speedUpVideoSuffix)

	normalVideosList := getFileList(dirname, normalVideoPrefix, normalVideoSuffix)
	speedUpVideo(dirname, normalVideosList, speedUpVideoName, speedUpMultiplier, normalVideoCleanAfter)

	// Upload to youtube
	description := fmt.Sprintf("Ускоренная запись вебкамеры со стройки ЖК \"Аист\"\nhttps://1strela.ru/aist/progress")
	privacy := "public"
	date := time.Now().Format("02 January 2006")
	title := fmt.Sprintf("Аист %s x%d", date, speedUpMultiplier)

	_, err := uploadVideo(speedUpVideoName, title, description, privacy)
	if err == nil && speedUpVideoCleanAfter {
		checkErr(os.Remove(speedUpVideoName))
	}
}

func getLastFile(dirname, prefix, suffix string) int {
	lastFile := 0
	dir, _ := ioutil.ReadDir(dirname)

	for _, file := range dir {
		filename := file.Name()

		if strings.HasPrefix(filename, prefix) && strings.HasSuffix(filename, suffix) {
			woPref := strings.TrimPrefix(filename, prefix)
			woPrefSuff := strings.TrimSuffix(woPref, suffix)
			number, err := strconv.Atoi(woPrefSuff)
			checkErr(err)

			if number > lastFile {
				lastFile = number
			}
		}
	}
	return lastFile
}

func getFileList(dirname, prefix, suffix string) []string {
	var fileList []string
	dir, _ := ioutil.ReadDir(dirname)

	for _, file := range dir {
		filename := file.Name()

		if strings.HasPrefix(filename, prefix) && strings.HasSuffix(filename, suffix) {
			filePath := fmt.Sprintf("%s%s", dirname, filename)
			fileList = append(fileList, filePath)
		}
	}
	return fileList
}

func checkErr(err error) {
	if err != nil {
		// panic(err)
		fmt.Println(err)
	}
}
