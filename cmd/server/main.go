package main

import (
	"log"
	"os"

	"github.com/eightlay/rummikub-server/iternal/server"
)

func main() {
	file, err := os.OpenFile("runtime.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(file)

	server.StartServer()
}
