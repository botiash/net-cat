package system

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type UserThread struct {
	Conn net.Conn
	Name string
}

type TotalUsers struct {
	Mu    *sync.Mutex
	Users map[string]net.Conn
}

var (
	totalUsers   = make(map[string]net.Conn, 10)
	InitialJoins bool
)

func CreateTotalUsers(mu *sync.Mutex) *TotalUsers {
	return (&TotalUsers{
		Mu:    mu,
		Users: make(map[string]net.Conn, 10), // returns an initialized map
	})
}

func (t *TotalUsers) AddUser(username string, user net.Conn) {
	t.Users[username] = user
}

func (t *TotalUsers) addUser(username string, user net.Conn) {
	t.Mu.Lock()
	t.AddUser(username, user)
	t.Mu.Unlock()
}

func CreateNewThread(conn net.Conn) *UserThread {
	return &UserThread{
		Conn: conn,
	}
}

func (user *UserThread) UserHandler(chBroad chan BroadCastStatus, chMsg chan BroadCastMessage, chat *Chat) {
	defer user.Conn.Close()

	intro, err := Intro() // Preparing the Intro message to the user
	if err != nil {
		if !Logger(err) {
			log.Println(err)
		}
		return // Defer closes connection
	}
	_, err = user.Conn.Write([]byte(intro)) // Writing the Intro message to the user
	if err != nil {
		if !Logger(err) {
			log.Println(err)
		}
		return // Defer closes connection
	}
	err = user.AddNewName(chat.Users) // adding a new user name
	if err != nil {
		if !Logger(err) {
			log.Println(err)
		}
		return // Defer closes connection
	}
	chat.Users.Mu.Lock()
	if len(chat.Users.Users) >= 10 { // Checking if the lobby is full
		user.LobbyIsFull()
		chat.Users.Mu.Unlock()
		return // Defer closes connection
	}
	chat.Users.Mu.Unlock()

	chat.Users.addUser(user.Name, user.Conn) // Adding the user to the totalUsers pool
	log.Printf("User %s joined the chat\n", user.Name)
	printErr := chat.PrintAllHistory(user) // Printing history to the user
	if printErr != nil {
		RemoveUser(chat, user.Name, chBroad)
		if !Logger(err) {
			log.Println("Error Printing History to user: ", user.Name, user.Conn)
		}
		return // Defer closes connection
	}
	chBroad <- BroadCastStatus{
		Name:        user.Name,
		IsConnected: true,
	}

	reader := bufio.NewReader(user.Conn) // Create user input reader

	for { // Listen for user input while connected
		_, err := fmt.Fprint(user.Conn, FormatMsg(user.Name, ""))
		if err != nil {
			if !Logger(err) {
				log.Printf("Error Printing Data to user %s\n", user.Name)
			}
			RemoveUser(chat, user.Name, chBroad)
			return // Defer closes connection
		}
		text, err := reader.ReadString('\n') // Store the msg that was sent by a user
		if err == io.EOF {
			log.Printf("User %s left the chat\n", user.Name)
			RemoveUser(chat, user.Name, chBroad)
			return // Defer closes connection
		}
		if err != nil {
			if !Logger(err) {
				log.Println(err)
			}
			RemoveUser(chat, user.Name, chBroad)
			return // Defer closes connection
		}
		trimmedText := strings.TrimSpace(text)
		if ValidText(trimmedText) { // Sending the valid msg to the channel
			chMsg <- BroadCastMessage{
				Name: user.Name,
				Msg:  FormatMsg(user.Name, trimmedText), // Sending the formatted text with dates and a name
				User: user,
			}
		}
	}
}

func (user *UserThread) AddNewName(users *TotalUsers) error {
	reader := bufio.NewReader(user.Conn)
	for { // Reading users input
		newName, err := reader.ReadString('\n')
		if err != nil { // If there is an error, remove user from the pool and log the error
			return err
		}
		newName = strings.TrimSpace(newName)

		valid, validErr := ValidName(user, newName, users) // Checking for name validity
		if valid == true && validErr == nil {
			user.Name = newName
			return nil
		}
		if validErr != nil { // If there is an error, return and remove user from the pool and log the error
			return validErr
		}
	}
}

func RemoveUser(chat *Chat, user string, chBroad chan BroadCastStatus) {
	chat.Users.Mu.Lock()
	delete(chat.Users.Users, user)
	chat.Users.Mu.Unlock()
	chBroad <- BroadCastStatus{IsConnected: false, Name: user}
}
