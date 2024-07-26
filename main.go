package main

import (
	"fmt"
	"ncrypt/encryptor"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Welcome to Ncrpyt")

	godotenv.Load()

	env := os.Getenv("PORT")

	fmt.Println(env)

	key := "12345"
	e := encryptor.Encrypt("Hello!123", key)
	d := encryptor.Decrypt(e, key)

	fmt.Println(e)
	fmt.Println(d)
}
