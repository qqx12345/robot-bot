package db

import (
	"log"
)

func Query() {
	client,err:=GlobalPool.Get()
	if err!=nil {
		log.Printf("连接获取失败：%v",err)
	}
	defer GlobalPool.Put(client)
	
}