package main

import (
	"github.com/joho/godotenv"
	"github.com/zercos/oauth-tower/internal/api"
)

func main() {
	godotenv.Load()
	serv := api.CreateServer()
	api.Run(serv)
}
