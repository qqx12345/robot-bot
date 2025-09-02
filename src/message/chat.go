package message

import (
	"net/http"
	"path"
	"encoding/json"
	"bytes"
	"log"
	"github.com/robot/src/token"
	"github.com/robot/src/db"
	"io"
)

func Chat(data map[string]interface{},ID string) {
	content:=data["content"].(string)
	id:=data["id"].(string)
	openid:=data["author"].(map[string]interface{})["user_openid"].(string)

    vector_usr:=Vector(content)

	content_llm:=Qwen(content)
	body:=map[string]interface{}{
		"content":content_llm,
		"msg_id":id,
		"msg_type":0,
	}

	vector_llm:=Vector(content_llm)

    insert_data:=map[string]interface{}{
		"text":[]string{content,content_llm},
		"open_id":[]string{openid,openid},
		"role":[]string{"user","AI"},
		"vector":[][]float32{vector_usr,vector_llm},
	}

    db.Insert(insert_data)

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