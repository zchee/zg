package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultHistorySize = 600
	fieldSep           = "\x00"
)

var (
	dataFile    string
	historySize int64
)

func getDataFileEnv() string {
	if dataFile = os.Getenv("_ZG_DATA"); len(dataFile) == 0 {
		dataFile = os.Getenv("HOME") + string(os.PathSeparator) + ".zg"
	}
	return dataFile
}

func getDataFile() (*os.File, error) {
	f := getDataFileEnv()
	_, err := os.Stat(f)
	if err != nil {
		fmt.Printf("%s: does not exist. Create %s? [y,n] :", f, f)
		if ask() {
			os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			return nil, err
		}
	}
	fd, _ := os.Open(f)
	defer fd.Close()

	return fd, nil
}

func getHistorySize() int64 {
	if historySize, _ = strconv.ParseInt(os.Getenv("_ZG_HISTORY_SIZE"), 10, 64); historySize < 1 {
		historySize = defaultHistorySize
	}
	return historySize
}
