package main

import (
	"fmt"
	"log"
	"os"

	"01.alem.school/git/aseitkha/net-cat/server"
	"01.alem.school/git/aseitkha/net-cat/system"
)

func main() {
	port := "8989"
	if len(os.Args) > 2 {
		log.Println("[USAGE]: ./TCPChat $port")
		return
	}

	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	addr := fmt.Sprintf("localhost:%s", port)
	serv, err := server.CreateNewServer(addr) // Starts the connection at the specified address
	if err != nil {                           // Checks for invalid port names etc.
		if !system.Logger(err) {
			log.Println(err)
		}
		return
	}
	defer serv.Listener.Close() // Defer closing the listener

	log.Println("Started the server at ", addr)
	if err := serv.RunServer(); err != nil { // Running the server
		if !system.Logger(err) {
			log.Println(err)
			return
		}
	}
}
