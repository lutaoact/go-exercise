package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type (
	pageInfo struct {
		Page     uint `json:"page,omitempty"`
		PerPage  uint `json:"perPage"`
		TotalNum uint `json:"totalNum"`
	}

	payload struct {
		Hello string `json:"hello"`
	}

	Response struct {
		Message  string    `json:"message"`
		PageInfo pageInfo  `json:"pageInfo"`
		SentAt   time.Time `json:"sentAt"`
		Payload  payload   `json:"payload"`
	}
)

func main() {
	jsonMarshal()
	//	jsonUnmarshal2()
	//	jsonUnmarshal()
	//	jsonUnmarshal2Map()
}

func decoder() {
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

var JSON = `{
	"message": "ok",
	"pageInfo": {
		"page": 1,
		"perPage": 10,
		"totalNum": 10
	},
	"sentAt": "1992-03-16T15:23:05+08:00",
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
	fmt.Printf("r.PageInfo.PerPage = %+v\n", r.PageInfo.PerPage)
	fmt.Printf("r.Payload.Hello = %+v\n", r.Payload.Hello)
	fmt.Printf("r = %+v\n", r)
}

var JSON2 = `
	{
		"perPage": 2,
		"totalNum": 3
	}
`

func jsonMarshal() {
	p := &pageInfo{
		PerPage:  2,
		TotalNum: 3,
	}
	//如果json tag中指明omitempty，则marshal为json时，不存在的字段会忽略掉
	//但这个字段不影响unmarshal，不管是否有omitempty，结构体对象都会初始化字段
	data, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(string(data))
}

func jsonUnmarshal2() {
	var p pageInfo
	err := json.Unmarshal([]byte(JSON2), &p)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	//	fmt.Printf("p.Page = %+v\n", p.Page)
	fmt.Printf("p.PerPage = %+v\n", p.PerPage)
	fmt.Printf("p.TotalNum = %+v\n", p.TotalNum)
	fmt.Printf("p = %+v\n", p)
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
