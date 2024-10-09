package main

import (
	"fmt"
	"log"
	"ncrypt/controllers"
	"ncrypt/utils"
	"ncrypt/utils/logger"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func deleteOldLogs() {
	path := "logs"

	//Always keep 5 most recent logs
	LOGS_TO_KEEP := 5

	// Get list of all files
	files, _ := os.ReadDir(path)

	// Sort files based on modified time in descending order
	sort.Slice(files, func(i, j int) bool {
		info1, _ := files[i].Info()
		info2, _ := files[j].Info()

		return info1.ModTime().After(info2.ModTime())
	})

	count := 0
	for _, file := range files {
		if count >= LOGS_TO_KEEP {
			// Delete files with .log type
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
				err := os.Remove(path + "/" + file.Name())
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}

		count += 1
	}
}

func main() {
	fmt.Println("Welcome to Ncrpyt")

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//Loading env
	godotenv.Load(".env")

	deleteOldLogs()

	utils.AssignDynamicPort()

	gin.DefaultWriter = logger.Log.Writer()

	//web server
	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	base_path := server.Group("")

	system_controller := new(controllers.SystemController)
	system_controller.Init()
	system_controller.RegisterRoutes(base_path)

	login_controller := new(controllers.LoginDataController)
	login_controller.Init()
	login_controller.RegisterRoutes(base_path)

	note_controller := new(controllers.NoteController)
	note_controller.Init()
	note_controller.RegisterRoutes(base_path)

	master_password_controller := new(controllers.MasterPasswordController)
	master_password_controller.Init()
	master_password_controller.RegisterRoutes(base_path)

	go func() {
		defer logger.Close()
	}()

	logger.Log.Printf("Starting server on %s", utils.PORT)
	server.Run(":" + utils.PORT)
}
