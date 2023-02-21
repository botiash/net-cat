package system

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func CreateLogger() {
	_, err := os.Stat("logger.txt")
	if os.IsNotExist(err) { // Checking if logger already exists
		file, err := os.Create("logger.txt")
		defer file.Close()
		if err != nil {
			log.Println("Failed to create a logger file")
			return
		}
	}
}

func Logger(err error) bool {
	successLog := false
	file, openErr := os.OpenFile("logger.txt", os.O_WRONLY, 0o644) // Open the file
	defer file.Close()
	if openErr != nil {
		log.Println("Failed to open the logger")
		return successLog
	}
	pc, filename, line, ok := runtime.Caller(1)
	if ok {
		logMsg := fmt.Sprintf("Error found in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err) // Setup the log
		_, err := fmt.Fprintln(file, logMsg)                                                                    // Write the log to the file
		if err != nil {
			log.Println("Failed to write to the logger")
			return successLog
		}
		successLog = true // Successfully wrote to the logger
	}
	return successLog
}
