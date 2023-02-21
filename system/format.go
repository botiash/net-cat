package system

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func Intro() (string, error) {
	fileText, err := ioutil.ReadFile("assets/logo.txt")
	if err != nil {
		log.Println("Couldn't open the welcome logo file")
		return "", err
	}
	introMsg := string(fileText) + "\n[ENTER YOUR NAME]:"
	return introMsg, nil
}

func FormatMsg(name string, msg string) string { // Formats the message properly
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	formattedInput := fmt.Sprintf("[%s][%s]:%s", formattedTime, name, msg)
	return strings.TrimSpace(formattedInput)
}

func ValidText(text string) bool { // Checks for the validity of the text
	return !(text == "" || (text[0] >= 0 && text[0] <= 31))
}

func ValidName(user *UserThread, name string, users *TotalUsers) (bool, error) { // Checks for the validity of user's name
	for _, ch := range name {
		if ch >= 0 && ch <= 31 {
			_, err := fmt.Fprint(user.Conn, "Unsopported name format\nPlease try another name...\n[ENTER YOUR NAME]:")
			if err != nil {
				return false, err // Checks if the user has been disconnected from the server
			}
			return false, nil
		}
	}
	if name == "" {
		_, err := fmt.Fprint(user.Conn, "Unsopported name format\nPlease try another name...\n[ENTER YOUR NAME]:")
		if err != nil {
			return false, err // Checks if the user has been disconnected from the server
		}
		return false, nil
	}
	users.Mu.Lock()
	if _, ok := users.Users[name]; ok {
		_, err := fmt.Fprint(user.Conn, "Specified name has already been taken\nPlease try another name...\n[ENTER YOUR NAME]:")
		if err != nil {
			users.Mu.Unlock()
			return false, err // Checks if the user has been disconnected from the server
		}
		users.Mu.Unlock()
		return false, nil
	}
	users.Mu.Unlock()
	return true, nil
}

func (user *UserThread) LobbyIsFull() { // Lobby is full msg
	_, err := fmt.Fprintf(user.Conn, "The lobby is currently at its full capacity...\nPlease try again later...\n")
	if err != nil {
		if !Logger(err) {
			log.Println(err)
		}
	}
}
