package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Welcome to Ncrpyt")

	godotenv.Load()

	env := os.Getenv("PORT")

	fmt.Println(env)
}
