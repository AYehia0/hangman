package main

import (
	"github.com/AYehia0/hangman/server"
	"github.com/charmbracelet/log"
)

func main() {

	host := "127.0.0.1"
	port := 1337 // love me ?

	log.Info("Running the SSH server...")
	server.Start(host, port)
}
