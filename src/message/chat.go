package message

import (
	"net/http"
	"path"
	"encoding/json"
	"bytes"
	"log"
	"github.com/robot/src/token"
)

func Chat(data map[string]interface{},ID string) {
	content:=data["content"].(string)
	id:=data["id"].(string)
	openid:=data["author"].(map[string]interface{})["user_openid"].(string)
	body:=map[string]interface{}{
		"content":content,
		"msg_id":id,
		"msg_type":0,
	}
	jsonValue, err := json.Marshal(body)
	if err != nil {
		log.Printf("JSON序列化失败: %v", err)
	}
    url := "https://api.sgroup.qq.com/v2/users/" + path.Join(openid, "messages")
	req,err:=http.NewRequest("POST",url,bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
	}
	req.Header.Set("Authorization",token.QQtoken.GetToken())
	req.Header.Set("Content-Type", "application/json")

	client:=&http.Client{}
	_, err = client.Do(req)
	if err != nil {
        log.Fatal("请求失败: %v", err)
    }
	

}