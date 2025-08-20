package src

import (
	"net/http"
	"os"
	"encoding/json"
	"log"
	"io"
	"bytes"
)

type token struct {
	token string
	endtime string
}

func Getoken () (string) {
	data:=map[string]string {
		"appId":os.Getenv("APPID"),
		"clientSecret":os.Getenv("BOTSECRET"),
	}
	jsonValue, _ := json.Marshal(data)
	req,err:=http.NewRequest("POST","https://bots.qq.com/app/getAppAccessToken",bytes.NewBuffer(jsonValue))
	if err != nil {
        log.Fatal(err)
    }
	client:=&http.Client{}
	res, err := client.Do(req)
	if err != nil {
        log.Fatal("请求失败:", err)
    }
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Printf("%s",string(body))
	return "1"
}