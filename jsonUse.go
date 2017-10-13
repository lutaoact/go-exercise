package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	pageInfo struct {
		Page     uint `json:"page"`
		PerPage  uint `json:"perPage"`
		TotalNum uint `json:"totalNum"`
	}

	payload struct {
		Hello string `json:"hello"`
	}

	Response struct {
		Message   string   `json:"message"`
		PageInfo  pageInfo `json:"pageInfo"`
		Timestamp uint64   `json:"timestamp"`
		Payload   payload  `json:"payload"`
	}
)

func main() {
	uri := "http://pdt.api.stockalert.cn/"
	resp, err := http.Get(uri)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()

	var r Response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(r.PageInfo.PerPage)
	fmt.Println(r.Payload.Hello)
	fmt.Println(r)

	jsonUnmarshal()
	jsonUnmarshal2Map()
}

var JSON = `{
	"message": "ok",
	"pageInfo": {
		"page": 1,
		"perPage": 10,
		"totalNum": 10
	},
	"timestamp": 1505125017201,
	"payload": {
		"hello": "girlfriend"
	}
}`

func jsonUnmarshal() {
	var r Response
	err := json.Unmarshal([]byte(JSON), &r)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(r.PageInfo.PerPage)
	fmt.Println(r.Payload.Hello)
	fmt.Println(r)
}

func jsonUnmarshal2Map() {
	var r map[string]interface{}
	err := json.Unmarshal([]byte(JSON), &r)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(r["message"])
	fmt.Println(r["pageInfo"].(map[string]interface{})["page"])
	fmt.Println(r["payload"].(map[string]interface{})["hello"])
}

/*
响应结构：
{
    "message": "ok",
    "pageInfo": {
        "page": 1,
        "perPage": 10,
        "totalNum": 10
    },
    "timestamp": 1505125017201,
    "payload": {
        "hello": "girlfriend"
    }
}
*/