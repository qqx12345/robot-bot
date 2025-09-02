package message

import (
	"net/http"
	"path"
	"encoding/json"
	"bytes"
	"log"
	"github.com/robot/src/token"
	"io"
)

func SendToQQ(content_llm string,id string, openid string) {
	body:=map[string]interface{}{
		"content":content_llm,
		"msg_id":id,
		"msg_type":0,
	}

	jsonValue, _ := json.Marshal(body)
    url := "https://api.sgroup.qq.com/v2/users/" + path.Join(openid, "messages")
	req,err:=http.NewRequest("POST",url,bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
	}
	req.Header.Set("Authorization","QQBot "+token.QQtoken.GetToken())
	req.Header.Set("Content-Type", "application/json")

	client:=&http.Client{}
	res, err := client.Do(req)
	resq,_:=io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
        log.Printf("请求失败: %v", err)
    }else {
		log.Printf("响应内容: %s", string(resq)) 
	}
}