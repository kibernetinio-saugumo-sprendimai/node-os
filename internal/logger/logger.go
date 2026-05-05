package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"nodeos/internal/identity"
)

const LogPath = "/var/lib/nodeos/nodeos.crypt.log"

type LogEntry struct {
	Timestamp string `json:"t"`
	Level     string `json:"l"`
	Message   string `json:"m"`
	Signature string `json:"sig"`
}

func Init() {
	_ = os.MkdirAll("/var/lib/nodeos", 0700)
}

func Log(level, msg string) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
	}

	// Sukuriame pasirašomą duomenų bloką (tik t, l, m)
	payload := fmt.Sprintf("%s|%s|%s", entry.Timestamp, entry.Level, entry.Message)
	
	// Pasirašome su mazgo privačiu raktu
	sig := identity.Sign([]byte(payload))
	entry.Signature = fmt.Sprintf("%x", sig)

	// Serializuojame į JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	// Rašome į failą (append)
	f, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(data)
	f.Write([]byte("\n"))
}
