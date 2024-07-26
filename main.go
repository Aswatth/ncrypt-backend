package main

import (
	"fmt"
	"log"
	"ncrypt/controllers"
	"ncrypt/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Welcome to Ncrpyt")
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//Loading env
	godotenv.Load()

	env := os.Getenv("PORT")

	//web server
	server := gin.Default()

	server.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "I am alive")
	})

	base_path := server.Group("")

	data_service := new(services.DataService).Init(os.Getenv("FILE_NAME"))
	data_controller := new(controllers.DataController)
	data_controller.Init(data_service)
	data_controller.RegisterRoutes(base_path)

	server.Run(":" + env)
}
