package alert

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type AlertHash [32]byte

func NewAlertHash(status, msg string) AlertHash {
	raw := fmt.Sprintf("%s:%s:%d", status, msg, time.Now().UnixNano())
	return sha256.Sum256([]byte(raw))
}

func (h AlertHash) String() string {
	return hex.EncodeToString(h[:])
}

func (h AlertHash) Prefix() string {
	return hex.EncodeToString(h[:4])
}
