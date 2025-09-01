package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/robot/src"
	"github.com/robot/src/db"
	"github.com/robot/src/message"
)

var GlobalPool *db.Pool

type Payload struct {
    ID string
    OP int
    D  map[string]interface{}
    S  int
    T  string
}

type HandlerFunc func(data map[string]interface{},ID string,T string) (interface{}, error)

var handlers = map[int]HandlerFunc{
	13: src.Sign,
    0: message.Message,
}

func app(writer http.ResponseWriter, request *http.Request) {
    httpBody, _ := io.ReadAll(request.Body)
    log.Printf("Chat data:\n%s", string(httpBody))
    defer request.Body.Close()
    payload := &Payload{}
    if err := json.Unmarshal(httpBody, payload); err != nil {
        http.Error(writer, "解析JSON失败", http.StatusBadRequest)
        return
    }
	res,err:= handlers[payload.OP](payload.D,payload.ID,payload.T)
    if err != nil {
        http.Error(writer, "中间件失败", http.StatusBadRequest)
        return
    }

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(res)
}

func factory() (*milvusclient.Client, error) {
    ctx := context.Background()
    client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
        Address:  os.Getenv("MILVUS_ADDRESS"),
    })
    if err != nil {
        return nil, err
    }
    return client, nil
}

func main() {
    var err error
    GlobalPool,err = db.Newpool(factory, 10)
    if err != nil {
        log.Printf("连接池初始化失败: %v", err)
    } else {
        log.Printf("连接池初始化成功")
    }
    
    // 在程序退出时关闭连接池
    defer func() {
        if GlobalPool != nil {
            GlobalPool.Close()
            log.Printf("连接池已关闭")
        }
    }()
    
	http.HandleFunc("/", app)
	http.ListenAndServe(":2345", nil)
}