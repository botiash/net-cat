package server

import (
	"net"
	"sync"

	"01.alem.school/git/aseitkha/net-cat/system"
)

type Server struct {
	Listener net.Listener
	address  string
}

func CreateNewServer(addr string) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Listener: listener,
		address:  addr,
	}, nil
}

func (server *Server) RunServer() error {
	chBroadCast := make(chan system.BroadCastStatus)     // Channel for User Status changes
	chMessage := make(chan system.BroadCastMessage)      // Channel for User Messages
	var muChat sync.Mutex                                // Mutex for Chat object
	var mu sync.Mutex                                    // Mutex for Total Users
	totalUsers := system.CreateTotalUsers(&mu)           // Creating a new total users struct
	chat := system.EstablishNewChat(&muChat, totalUsers) // Establishing a new Chat object
	system.CreateLogger()                                // Logger.txt creation
	go chat.BroadCastRoutine(chBroadCast, chMessage)     // Creating a Broadcaster in a separate thread
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			return err
		}

		newUserThread := system.CreateNewThread(conn) // Creating a user struct with net.Conn

		go newUserThread.UserHandler(chBroadCast, chMessage, chat) // User Handler
	}
}
