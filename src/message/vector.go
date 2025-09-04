package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Vector(data string) []float32{
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	url := "https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings"

	jsonData := []byte(fmt.Sprintf(`{
		"model": "text-embedding-v4",
		"input": "%s",
		"dimensions": "512",
		"encoding_format": "float"
	}`, data))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	json.Unmarshal(body, &result)

	if len(result.Data) == 0 {
		log.Println("暂无向量")
		return nil
	}
	embedding := result.Data[0].Embedding
	return embedding
}
