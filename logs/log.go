//Usage:
//import "logs"
//
//logs.Info("info")
//logs.Warn("warning")
//logs.Debug("debug")
//logs.Error("error")
//logs.Critical("critical")
package logs

import (
	"fmt"
	"log"
)

// output INFO leve message, should have Format param
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I]"+format, v...)
	log.Println(msg)
}

// simple output INFO leve message
func I(v ...interface{}) {
	Info("%s", v...)
}

// output WARN leve message, should have Format param
func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[W]"+format, v...)
	log.Println(msg)
}

// simple output WARN leve message
func W(v ...interface{}) {
	Warn("%s", v...)
}

// output DEBUG leve message, should have Format param
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D]"+format, v...)
	log.Println(msg)
}

// simple output DEBUG leve message
func D(v ...interface{}) {
	Debug("%s", v...)
}

// output ERROR leve message, should have Format param
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E]"+format, v...)
	log.Println(msg)
}

// simple output ERROR leve message
func E(v ...interface{}) {
	Error("%s", v...)
}

// output CRITICAL leve message, should have Format param
func Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf("[C]"+format, v...)
	log.Println(msg)
}

// simple output CRITICAL leve message
func C(v ...interface{}) {
	Critical("%s", v...)
}
