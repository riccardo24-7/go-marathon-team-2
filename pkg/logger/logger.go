package logger

import (
	"log"
	"os"
	"sync"
)

var LogChan = make(chan *message)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	DEBUG LogLevel = "DEBUG"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

type message struct {
	err error
	lvl LogLevel
	msg string
}

func NewMessage(level LogLevel, msgInfo string, err error) *message {
	return &message{
		msg: msgInfo,
		lvl: level,
		err: err,
	}
}

func NewLogger(file *os.File, debug bool, wg *sync.WaitGroup) {

	log.SetOutput(file)
	defer close(LogChan)

	for {
		select {
		case msg := <-LogChan:
			wg.Add(1)
			if debug {
				log.Print(msg.lvl+" | ", msg.msg, "\n", "error: ", msg.err)
			} else {
				if msg.lvl != DEBUG {
					log.Print(msg.lvl+" | ", msg.msg, "\n", "error: ", msg.err)
				}
			}
			wg.Done()
		}
	}
}

func LogMessageInfo(msg string, err error) {
	LogChan <- NewMessage(INFO, msg, err)
}

func LogMessageError(msg string, err error) {
	LogChan <- NewMessage(ERROR, msg, err)
}

func LogMessageFatal(msg string, err error) {
	LogChan <- NewMessage(FATAL, msg, err)
	log.Fatalln(err)
}
