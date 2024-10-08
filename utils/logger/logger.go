package logger

import (
	"log"
	"os"
	"strings"
	"time"
)

var (
	Log  *log.Logger
	file *os.File
)

func init() {
	dir := "logs"
	os.Mkdir(dir, os.ModePerm)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic("ERROR creating logs directory")
		}
	}

	logpath := dir + "\\log-" + time.Now().Format(time.RFC3339) + ".log"
	logpath = strings.ReplaceAll(logpath, ":", "-")

	//    flag.Parse()
	file, err := os.Create(logpath)

	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + logpath)
}

func Close() {
	file.Close()
}
