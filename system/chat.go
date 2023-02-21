package system

import (
	"fmt"
	"log"
	"sync"
)

type Chat struct {
	HistoryBuffer []string
	mu            *sync.Mutex
	Users         *TotalUsers
}

type BroadCastMessage struct {
	Name string
	Msg  string
	User *UserThread
}

type BroadCastStatus struct {
	IsConnected bool
	Name        string
}

func EstablishNewChat(mu *sync.Mutex, users *TotalUsers) *Chat {
	return &Chat{
		mu:    mu,
		Users: users,
	}
}

func (chat *Chat) BroadCastRoutine(chBroadCast chan BroadCastStatus, chMessage chan BroadCastMessage) {
	for {
		select {
		case status := <-chBroadCast: // Status BroadCast
			message := fmt.Sprintf("%s has left the channel", status.Name)
			if status.IsConnected {
				message = fmt.Sprintf("%s has joined the channel", status.Name)
			}
			chat.send(message, status.Name, chBroadCast)

		case incomingMsg := <-chMessage: // Message BroadCast
			chat.send(incomingMsg.Msg, incomingMsg.Name, chBroadCast)
		}
	}
}

func (chat *Chat) send(msg string, name string, chBroad chan BroadCastStatus) {
	chat.mu.Lock()
	chat.HistoryBuffer = append(chat.HistoryBuffer, msg)
	chat.mu.Unlock()

	chat.Users.Mu.Lock()
	defer chat.Users.Mu.Unlock()
	for user, Conn := range chat.Users.Users { // Range over all users and print msg to others
		if user != name {
			_, err := fmt.Fprint(Conn, "\n"+msg+"\n"+FormatMsg(user, ""))
			if err != nil {
				delete(chat.Users.Users, user) // Remove the user to whom u couldn't write since it means that he is disconnected
				chBroad <- BroadCastStatus{IsConnected: false, Name: user}
				if !Logger(err) {
					log.Println("Error Printing Data to user: ", user, Conn)
				}
				continue
			}
		}
	}
}

func (chat *Chat) PrintAllHistory(user *UserThread) error {
	chat.mu.Lock()
	defer chat.mu.Unlock()
	for _, text := range chat.HistoryBuffer {
		_, err := fmt.Fprint(user.Conn, text+"\n")
		if err != nil {
			return err
		}
	}
	return nil
}
