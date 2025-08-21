package src

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"log"
	"strings"
	"os"
)

func Sign(data map[string]interface{},ID string,T string) (interface{}, error) {

	plain_token:= data["plain_token"].(string)
	event_ts:= data["event_ts"].(string)
	seed := os.Getenv("BOTSECRET")
	for len(seed) < ed25519.SeedSize {
		seed = strings.Repeat(seed, 2)
	}
	seed = seed[:ed25519.SeedSize]
	_, privateKey, err := ed25519.GenerateKey(strings.NewReader(seed))
	if err != nil {
		log.Println("ed25519 generate key failed:", err)
		return nil, err
	}
	var msg bytes.Buffer
	msg.WriteString(event_ts)
	msg.WriteString(plain_token)
	signature := hex.EncodeToString(ed25519.Sign(privateKey, msg.Bytes()))

	res := map[string]string{
		"plain_token":plain_token,
		"signature":signature,
	}
	return res, nil
}
