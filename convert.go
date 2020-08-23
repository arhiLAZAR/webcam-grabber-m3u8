package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	tmpVideoPrefix       = "tmp_video_"
	tmpVideoSuffix       = ".mp4"
	tmpUpVideoCleanAfter = true
)

func speedUpVideo(dirname string, normalVideosList []string, speedUpVideoName string, speedUpMultiplier int64, clearAfterNormal ...bool) {
	var tmpVideosList []string

	ptsMultilier := 1 / float64(speedUpMultiplier)
	setpst := fmt.Sprintf("setpts=%f*PTS", ptsMultilier)

	for _, normalVideoName := range normalVideosList {
		tmpVideoCount := getLastFile(dirname, tmpVideoPrefix, tmpVideoSuffix)
		tmpVideoCount++
		tmpVideoName := fmt.Sprintf("%s%s%06.f%s", dirname, tmpVideoPrefix, float64(tmpVideoCount), tmpVideoSuffix)

		cmd := exec.Command("ffmpeg", "-i", normalVideoName, "-r", "60", "-filter:v", setpst, tmpVideoName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err == nil {
			tmpVideosList = append(tmpVideosList, tmpVideoName)

			if clearAfterNormal[0] {
				checkErr(os.Remove(normalVideoName))
			}

		} else {
			checkErr(err)
		}

	}

	concatVideo(tmpVideosList, speedUpVideoName)
}

func concatVideo(tmpVideosList []string, speedUpVideoName string) {
	inputFilesList := dirname + "input_files.txt"

	makeInpitFilesList(tmpVideosList, inputFilesList)

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", inputFilesList, speedUpVideoName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err == nil && tmpUpVideoCleanAfter {
		for _, tmpVideo := range tmpVideosList {
			checkErr(os.Remove(tmpVideo))
		}
	} else {
		checkErr(err)
	}

	checkErr(os.Remove(inputFilesList))
}

func makeInpitFilesList(videosList []string, inputFilesList string) {

	file, err := os.OpenFile(inputFilesList, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	checkErr(err)
	defer file.Close()

	for _, video := range videosList {
		line := fmt.Sprintf("file %s\n", video)
		_, err := file.WriteString(line)
		checkErr(err)
	}

}
