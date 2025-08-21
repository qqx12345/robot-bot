package main

import (
	"encoding/json"
	"io"
	"net/http"
	"github.com/robot/src"
    "github.com/robot/src/message"
    "log"
)

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

func main() {
	http.HandleFunc("/",app)
	http.ListenAndServe(":2345", nil)
}