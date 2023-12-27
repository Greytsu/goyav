package spotify

import (
	"fmt"
	"time"
)

// Parses a string to Time using this layout : "2006-01-02T15:04:05.999999Z07:00"
func ParseExpiry(expiryStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.999999Z07:00"

	expiry, err := time.Parse(layout, expiryStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Now(), err
	}

	return expiry, nil
}
