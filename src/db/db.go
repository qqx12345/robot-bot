package db

import (
	"context"
	"fmt"
	"log"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type Querydata struct {
	Role, Text string
}

func Query(openID string, queryVector []float32) []Querydata {
	// 1. 获取连接并检查错误
	client, err := GlobalPool.Get()
	if err != nil {
		log.Printf("连接获取失败：%v", err)
		return []Querydata{}
	}
	defer GlobalPool.Put(client)

	// 2. 检查查询向量是否为空
	if len(queryVector) == 0 {
		log.Println("警告：查询向量为空")
		return []Querydata{}
	}

	ctx := context.Background()
	resultSets, err := client.Search(ctx, milvusclient.NewSearchOption(
		"demo_collection",
		5, // 最多返回5条结果
		[]entity.Vector{entity.FloatVector(queryVector)},
	).WithConsistencyLevel(entity.ClStrong).
		WithANNSField("text_dense_vector").
		WithFilter(fmt.Sprintf("user_openid == \"%s\"", openID)).
		WithOutputFields("text", "role"))

	// 3. 检查搜索结果错误
	if err != nil {
		log.Printf("搜索失败：%v", err)
		return []Querydata{}
	}

	// 4. 检查结果集是否为空
	if resultSets == nil {
		log.Println("提示：未找到相关结果")
		return []Querydata{}
	}

	back := []Querydata{}

	// 5. 安全地处理每一条结果
	for _, resultSet := range resultSets {
		// 检查结果集是否包含必要的列
		textColumn := resultSet.GetColumn("text")
		roleColumn := resultSet.GetColumn("role")

		if textColumn == nil || roleColumn == nil {
			log.Println("警告：结果集中缺少必要的列")
			continue
		}

		// 安全地获取字段数据
		textFieldData := textColumn.FieldData()
		roleFieldData := roleColumn.FieldData()

		if textFieldData == nil || roleFieldData == nil {
			log.Println("警告：无法获取字段数据")
			continue
		}

		// 安全地获取标量值
		textScalars := textFieldData.GetScalars()
		roleScalars := roleFieldData.GetScalars()

		if textScalars == nil || roleScalars == nil {
			log.Println("警告：无法获取标量值")
			continue
		}

		fmt.Println("text: ", textScalars)
		fmt.Println("role: ", roleScalars)
	}

	return back
}

func Insert(data map[string]interface{}) {
	// 获取连接
	client, err := GlobalPool.Get()
	if err != nil {
		log.Printf("连接获取失败：%v", err)
		return
	}
	defer GlobalPool.Put(client)

	// 检查数据是否完整
	requiredKeys := []string{"open_id", "role", "text", "vector"}
	for _, key := range requiredKeys {
		if _, exists := data[key]; !exists {
			log.Printf("警告：插入数据缺少必要的字段：%s", key)
			return
		}
	}

	ctx := context.Background()
	_, err = client.Insert(ctx, milvusclient.NewColumnBasedInsertOption("demo_collection").
		WithVarcharColumn("user_openid", data["open_id"].([]string)).
		WithVarcharColumn("role", data["role"].([]string)).
		WithVarcharColumn("text", data["text"].([]string)).
		WithFloatVectorColumn("text_dense_vector", 512, data["vector"].([][]float32)),
	)

	if err != nil {
		log.Printf("插入失败：%v", err)
	} else {
		log.Println("插入成功")
	}
}