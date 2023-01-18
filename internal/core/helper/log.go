package helper

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func InitializeLog() {
	logDir := Config.LogDir
	_ = os.Mkdir(logDir, os.ModePerm)

	f, err := os.OpenFile(logDir+Config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		LogEvent("ERROR", "creating log file"+err.Error())
		log.Fatalf("error opening file: %v", err)
	}
	log.SetFlags(0)
	log.SetOutput(f)
}

// level can be INFO or ERROR
func LogEvent(level string, message interface{}) {

	data, err := json.Marshal(struct {
		TimeStamp string      `json:"@timestamp"`
		Level     string      `json:"level"`
		AppName   string      `json:"app_name"`
		Message   interface{} `json:"message"`
	}{
		TimeStamp: time.Now().Format(time.RFC3339),
		AppName:   Config.ServiceName,
		Message:   message,
		Level:     level,
	})

	if err != nil {
		log.Fatal("Logevent helper: " + err.Error())
	}
	log.Printf("%s\n", data)

}
