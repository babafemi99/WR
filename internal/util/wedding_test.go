package util

import (
	"log"
	"testing"
)

func TestGenerateSpecialKey(t *testing.T) {
	log.Println(GenerateSpecialKey("A3840hf"))
}
