package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"
)

// Helper: return a sql.NullString
func NullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: s}
}

func IsAllowedImageExtension(ext string) bool {
	ext = strings.ToLower(ext)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	return allowed[ext]
}

// generateRandomString generates a cryptographically secure random string.
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		log.Printf("WARNING: Failed to generate random string: %v. Falling back to less secure method.", err)
		// Fallback for demonstration, but for real passwords, you should fatal or re-attempt.
		return fmt.Sprintf("%x", time.Now().UnixNano())[:length]
	}
	return base64.URLEncoding.EncodeToString(b)[:length] // Ensure it's exactly 'length' chars if base64 padding is an issue
}

// generateRandomDateOfBirth generates a random date of birth
// representing someone between 20 and 60 years old.
func GenerateRandomDateOfBirth() string {
	now := time.Now()
	minAgeYears := 18
	maxAgeYears := 80

	minDOB := now.AddDate(-maxAgeYears, 0, 0)
	maxDOB := now.AddDate(-minAgeYears, 0, 0)

	duration := maxDOB.Sub(minDOB)
	randomNanos := RandInt64(0, duration.Nanoseconds())
	randomDuration := time.Duration(randomNanos)

	randomDOBTime := minDOB.Add(randomDuration)

	// We only care about the date part, so we can set time to UTC midnight
	// and ignore the original time component of randomDOBTime.
	randomDOBOnly := time.Date(randomDOBTime.Year(), randomDOBTime.Month(), randomDOBTime.Day(),
		0, 0, 0, 0, time.UTC)

	// Format using "2006-01-02" which represents YYYY-MM-DD
	return randomDOBOnly.Format("2006-01-02") // <--- CHANGE IS HERE
}

// randInt64 generates a random int64 between min and max (inclusive).
// Note: For real applications, consider using a more robust random number generator
// and handling potential crypto/rand errors.
func RandInt64(min, max int64) int64 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	// Generate random bytes
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		log.Printf("ERROR: crypto/rand.Read failed: %v, falling back to math/rand (less secure)", err)
		// Fallback to less secure math/rand if crypto/rand fails (shouldn't happen often)
		// Ensure math/rand is seeded if used as a fallback
		return min + (time.Now().UnixNano() % (max - min + 1))
	}
	// Convert bytes to uint64
	val := uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16 | uint64(buf[3])<<24 |
		uint64(buf[4])<<32 | uint64(buf[5])<<40 | uint64(buf[6])<<48 | uint64(buf[7])<<56
	return min + int64(val%(uint64(max-min)+1))
}
