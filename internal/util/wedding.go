package util

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GenerateSpecialKey(ID string) string {
	// Convert the ID string to bytes
	idBytes := []byte(ID)

	// Create a new SHA-256 hash instance
	hash := sha256.New()

	// Write the ID bytes to the hash
	_, err := hash.Write(idBytes)
	if err != nil {
		// Handle error, if any
		return ""
	}

	// Sum the hash and get the hash value as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the first 4 bytes of the hash value to a hexadecimal string
	hashHex := hex.EncodeToString(hashBytes[:4])

	// Insert '-' after every 3 characters
	formattedKey := insertDashEveryN(hashHex, 3)

	return formattedKey
}

// insertDashEveryN inserts a '-' after every n characters in a string
func insertDashEveryN(str string, n int) string {
	var parts []string
	for len(str) > n {
		parts = append(parts, str[:n])
		str = str[n:]
	}
	parts = append(parts, str)
	return strings.Join(parts, "-")
}
