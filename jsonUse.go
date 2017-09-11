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
  "timestamp": 1505123768800,
  "payload": {
    "hello": "girlfriend"
  }
}
*/
