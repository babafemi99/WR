package util

import (
	"log"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	token, _, err := GenerateToken("12345", "super", "a@a.com", "7765", time.Now().Add(time.Hour*24*30))
	if err != nil {
		return
	}
	log.Println(token)
}
