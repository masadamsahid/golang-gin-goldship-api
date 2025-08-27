package helpers

import (
	"math/rand"
	"time"
)

func randomString(length int) string {
	// Seed the random number generator using the current time for different results each run.
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func GenerateTrackingNumber() string {
	currDate := time.Now()

	return "GS-" + currDate.Format("20060102150405") + randomString(8)
}
