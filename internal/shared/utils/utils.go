package utils

import (
	"crypto/rand"
	"math/big"
	"time"
)

func GenerateBookingCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	return "GLP-" + string(b)
}

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}
