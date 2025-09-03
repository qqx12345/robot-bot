package message

import (
	"github.com/robot/src/db"
)

func Chat(data map[string]interface{},ID string) {
	content:=data["content"].(string)
	id:=data["id"].(string)
	openid:=data["author"].(map[string]interface{})["user_openid"].(string)

    vector_usr:=Vector(content)

	back_info:=db.Query(openid,vector_usr)

	content_llm:=Qwen(content,back_info)

	go func(){
		SendToQQ(content_llm,id,openid)
	}()

	vector_llm:=Vector(content_llm)

    insert_data:=map[string]interface{}{
		"text":[]string{content,content_llm},
		"open_id":[]string{openid,openid},
		"role":[]string{"user","AI"},
		"vector":[][]float32{vector_usr,vector_llm},
	}

    db.Insert(insert_data)
}