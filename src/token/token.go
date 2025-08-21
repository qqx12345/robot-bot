package token

import (
	"net/http"
	"os"
	"encoding/json"
	"log"
	"io"
	"bytes"
	"time"
	"sync"
	"strconv"
)

type Token struct {
	token string
	endtime time.Time
    mu sync.RWMutex
}

var QQtoken = &Token{}

func (tm *Token) GetToken() (string) {
    // 第一次检查：快速路径（读锁）
    tm.mu.RLock()
    if tm.isValid() {
        token := tm.token
        tm.mu.RUnlock()
        return token
    }
    tm.mu.RUnlock()

    // 慢速路径：获取写锁
    tm.mu.Lock()
    defer tm.mu.Unlock()

    // 第二次检查：避免重复刷新
    if tm.isValid() {
        return tm.token
    }

    // 真正刷新token
    tm.token,tm.endtime= tm.refresh()
	return tm.token
}

func (tm *Token) isValid() bool {
	return tm.token != "" && time.Now().Before(tm.endtime)
}

func (tm *Token) refresh () (string, time.Time) {
	data:=map[string]string {
		"appId":os.Getenv("APPID"),
		"clientSecret":os.Getenv("BOTSECRET"),
	}
	jsonValue, _ := json.Marshal(data)
	req,err:=http.NewRequest("POST","https://bots.qq.com/app/getAppAccessToken",bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")

	if err != nil {
        log.Fatal(err)
    }
	client:=&http.Client{}
	curtime:=time.Now()
	res, err := client.Do(req)
	if err != nil {
        log.Fatal("请求失败:", err)
    }

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var result map[string]string
	json.Unmarshal(body, &result)
	accessToken := result["access_token"]
	num, err := strconv.Atoi(result["expires_in"])
	if err!=nil {
		log.Printf("转换错误: %v", err)
	}
	endtime := curtime.Add(time.Duration(num-60)*time.Second)

	return accessToken,endtime
}