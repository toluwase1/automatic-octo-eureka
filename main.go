package main

import (
	"github.com/joho/godotenv"
	"log"
	"wallet-engine/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if  err != nil {
		log.Printf("Error loading .env %s", err)
	}
	server.Start()
}
