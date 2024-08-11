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
	os.Mkdir("logs", os.ModePerm)
	logpath := "logs\\log-" + time.Now().Format(time.RFC3339) + ".log"
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
