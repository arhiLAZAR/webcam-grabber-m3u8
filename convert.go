package main

import (
	"fmt"
	"os"
	"os/exec"
)

func speedUpVideo(dirname string, normalVideosList []string, speedUpVideoName string, speedUpMultiplier int64, clearAfter ...bool) {
	inputFilesList := dirname + "input_files.txt"

	makeInpitFilesList(normalVideosList, inputFilesList)

	ptsMultilier := 1 / float64(speedUpMultiplier)
	setpst := fmt.Sprintf("setpts=%f*PTS", ptsMultilier)

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", inputFilesList, "-r", "60", "-filter:v", setpst, speedUpVideoName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err == nil && clearAfter[0] {
		for _, normalVideo := range normalVideosList {
			checkErr(os.Remove(normalVideo))
		}
	} else {
		checkErr(err)
	}

	checkErr(os.Remove(inputFilesList))
}

func makeInpitFilesList(normalVideosList []string, inputFilesList string) {

	file, err := os.OpenFile(inputFilesList, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	checkErr(err)
	defer file.Close()

	for _, normalVideo := range normalVideosList {
		line := fmt.Sprintf("file %s\n", normalVideo)
		_, err := file.WriteString(line)
		checkErr(err)
	}

}
