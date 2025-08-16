package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	apiKey := "" // 替换为你的API Key
	url := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	
	requestBody := map[string]interface{}{
		"model": "qwen-turbo", // 模型名称
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": "你好，介绍一下你自己",
				},
			},
		},
		"parameters": map[string]interface{}{
			"result_format": "message", // 返回格式
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}