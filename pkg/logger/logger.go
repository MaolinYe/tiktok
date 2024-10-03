package logger

import "log"

func InitLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("LOG: ")
}
