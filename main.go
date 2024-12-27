package main

import (
	"github.com/zercos/oauth-tower/internal/server"
)

func main() {
	serv := server.CreateServer()
	server.Run(serv)
}
