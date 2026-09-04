package logger

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"nodeos/internal/identity"
)

const LogPath = "/var/lib/nodeos/nodeos.crypt.log"

type LogEntry struct {
	Timestamp    string `json:"t"`
	Level        string `json:"l"`
	Message      string `json:"m"`
	PreviousHash string `json:"prev"`
	EntryHash    string `json:"hash"`
	Signature    string `json:"sig"`
}

func Init() {
	_ = os.MkdirAll("/var/lib/nodeos", 0700)
}

func lastEntryHash() (string, error) {
	f, err := os.Open(LogPath)
	if errors.Is(err, os.ErrNotExist) {
		return strings.Repeat("0", 64), nil
	}
	if err != nil {
		return "", err
	}
	defer f.Close()

	last := strings.Repeat("0", 64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil || entry.EntryHash == "" {
			return "", errors.New("audit log chain is malformed")
		}
		last = entry.EntryHash
	}
	return last, scanner.Err()
}

func Log(level, msg string) error {
	previousHash, err := lastEntryHash()
	if err != nil {
		return err
	}
	entry := LogEntry{
		Timestamp:    time.Now().UTC().Format(time.RFC3339Nano),
		Level:        level,
		Message:      msg,
		PreviousHash: previousHash,
	}

	// The signature covers the complete chain link, not only the message.
	payload := fmt.Sprintf("%s|%s|%s|%s", entry.Timestamp, entry.Level, entry.Message, entry.PreviousHash)
	digest := sha256.Sum256([]byte(payload))
	entry.EntryHash = fmt.Sprintf("%x", digest)
	
	sig, err := identity.Sign([]byte(payload))
	if err != nil {
		return err
	}
	entry.Signature = fmt.Sprintf("%x", sig)

	// Serializuojame į JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	// Rašome į failą (append)
	f, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return f.Sync()
}
