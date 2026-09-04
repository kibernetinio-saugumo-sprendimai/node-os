package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nodeos/internal/config"
	"nodeos/internal/logger"
)

// telegramMessage – Telegram API payload
type telegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// Info – siunčia informacinį pranešimą
func Info(msg string) {
	if canSend("INFO") {
		updateMemory("INFO", msg)
		sendFormattedAlert("INFO", "ℹ️", msg)
	}
}

// Success – siunčia sėkmės pranešimą
func Success(msg string) {
	if canSend("SUCCESS") {
		updateMemory("SUCCESS", msg)
		sendFormattedAlert("SUCCESS", "✅", msg)
	}
}

// Warn – siunčia įspėjimą
func Warn(msg string) {
	if canSend("WARN") {
		updateMemory("WARN", msg)
		sendFormattedAlert("WARN", "⚠️", msg)
	}
}

// Critical – siunčia kritinį pranešimą
func Critical(msg string) {
	if canSend("CRITICAL") {
		updateMemory("CRITICAL", msg)
		sendFormattedAlert("CRITICAL", "🚨", msg)
	}
}

// Destruct – siunčia pranešimą apie sunaikinimą
func Destruct(msg string) {
	if canSend("DESTRUCT") {
		updateMemory("DESTRUCT", msg)
		sendFormattedAlert("DESTRUCT", "💥", msg)
	}
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// sendFormattedAlert – siunčia suformatuotą alertą
func sendFormattedAlert(status string, emoji string, msg string) error {
	cfg := config.LoadConfig()

	// Log to crypto-log (Signed)
	if err := logger.Log(status, msg); err != nil {
		return err
	}

	// Update internal engines
	UpdateAwareness(status)
	updatePersonality()

	// Jei nėra Telegram konfigo – tyliai išeinam
	if cfg.TelegramToken == "" || cfg.TelegramChatID == "" {
		return nil
	}

	// Sukuriam per-alert hash (be global state)
	h := NewAlertHash(status, msg)
	ref := h.Prefix()

	// Pagrindinė žinutė
	fullMsg := emoji + " [" + status + "] " + msg + " | #" + ref

	// Awareness + Personality sluoksniai
	fullMsg = fullMsg + AwarenessComment()
	fullMsg = fullMsg + personalityComment()

	// Pilnas hash auditui
	fullMsg = fullMsg + "\n#hash: " + h.String()

	// Telegram API
	url := "https://api.telegram.org/bot" + cfg.TelegramToken + "/sendMessage"

	body := telegramMessage{
		ChatID: cfg.TelegramChatID,
		Text:   fullMsg,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram API returned %s", resp.Status)
	}
	return nil
}
