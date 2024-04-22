package util

import (
	"log"
	"testing"
)

func TestGenerateSpecialKey(t *testing.T) {
	log.Println(GenerateSpecialKey("A3840hf"))
}

func TestEncodeCID(t *testing.T) {
	log.Println(EncodeCID("223345566"))
}

func TestDecodeCID(t *testing.T) {
	cid, err := DecodeCID("NTA1ODk1ODAxM3x8Sm9obiBhbmQgSmFuZQ==")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	log.Println(cid)
}
