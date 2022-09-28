package log

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func NewLogger() *log.Logger {
	// Open Log file
	logfile, err := os.OpenFile(fmt.Sprint(viper.Get("LOGFILE")), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	// Create error logger
	logger := log.New(logfile, "[ERROR] ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)

	return logger
}
