package db

import (
	"log"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"context"
	"fmt"
)

type Querydata struct {Role, Text string}

func Query(open_id string,queryVector []float32) []Querydata {
	client,err:=GlobalPool.Get()
	if err!=nil {
		log.Printf("连接获取失败：%v",err)
	}
	defer GlobalPool.Put(client)
	ctx:=context.Background()
	resultSets, err := client.Search(ctx, milvusclient.NewSearchOption(
        "demo_collection",
        5,
        []entity.Vector{entity.FloatVector(queryVector)},
    ).WithConsistencyLevel(entity.ClStrong).
        WithANNSField("text_dense_vector").
        WithFilter(fmt.Sprintf("user_openid == \"%s\"", open_id)))
	if err != nil {
        fmt.Println(err.Error())
    }
	back := []Querydata{}
	for _, resultSet := range resultSets {
		text:=resultSet.GetColumn("text").FieldData().GetScalars()
		role:=resultSet.GetColumn("role").FieldData().GetScalars()
		fmt.Println("text: ", text)
		fmt.Println("role: ", role)
    }
	return back
}

func Insert(data map[string]interface{}) {
	client,err:=GlobalPool.Get()
	if err!=nil {
		log.Printf("连接获取失败：%v",err)
	}
	defer GlobalPool.Put(client)
	ctx:=context.Background()
    _, err = client.Insert(ctx, milvusclient.NewColumnBasedInsertOption("demo_collection").
    WithVarcharColumn("user_openid", data["open_id"].([]string)).
    WithVarcharColumn("role", data["role"].([]string)).
    WithVarcharColumn("text", data["text"].([]string)).
    WithFloatVectorColumn("text_dense_vector", 512, data["vector"].([][]float32)),
)
	if err != nil {
        fmt.Println(err.Error())
    } else {
		fmt.Println("插入成功")
	}
}