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

	port := os.Getenv("PORT")

	//web server
	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	base_path := server.Group("")

	login_service := new(services.LoginService)
	login_service.Init()
	login_controller := new(controllers.LoginController)
	login_controller.Init(login_service)
	login_controller.RegisterRoutes(base_path)

	master_password_service := new(services.MasterPasswordService)
	master_password_service.Init()
	master_password_controller := new(controllers.MasterPasswordController)
	master_password_controller.Init(master_password_service)
	master_password_controller.RegisterRoutes(base_path)

	server.Run(":" + port)
}
