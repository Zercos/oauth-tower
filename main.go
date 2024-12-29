package main

import (
	"github.com/zercos/oauth-tower/internal/api"
)

func main() {
	serv := api.CreateServer()
	api.Run(serv)
}
