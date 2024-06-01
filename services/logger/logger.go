package logger

import (
	"fmt"
	"log"
)

type LogType uint8

const (
	INFO = iota
	WARNING
	ERROR
)

func Log(logType LogType, msg string, v... any) {
	finalMsg := "["
	switch logType {
	case INFO:
		finalMsg += "INFO"
		break;
	case WARNING:
		finalMsg += "WARNING"
		break;
	case ERROR:
		finalMsg += "ERROR"
		break;
	}
	finalMsg += "] " + fmt.Sprintf(msg, v...)

	log.Println(finalMsg)
}
