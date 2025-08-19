package sign

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"log"
	"strings"
	"os"
)

func getStringField(data map[string]interface{}, key string) string {
	if val, exists := data[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return "" // 返回空字符串作为默认值
}

func Sign(data map[string]interface{}) (interface{}, error) {

	plain_token:= getStringField(data, "plain_token")
	event_ts:= getStringField(data, "event_ts")
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
