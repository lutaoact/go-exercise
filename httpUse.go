package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Header map[string][]string

type Person struct {
	Name     string    `form:"name"     json:"name,omitempty"                                           binding:"required"`
	Address  string    `form:"address"  json:"address,omitempty"                                        binding:"required"`
	Birthday time.Time `form:"birthday" json:"birthday,omitempty" time_format:"2006-01-02" time_utc:"1" binding:"required"`
}

func main() {
	//	TestHttpPost()
	TestHttpGet()
}

func TestHeader() {
	header := Header{
		"Content-Type": []string{"application/json"},
	}
	for k, v := range header {
		fmt.Println(k, strings.Join(v, ""))
	}
}

func TestHttpGet() {
	resp, err := http.Get("http://xxx.lutaoact.com")
	fmt.Printf("resp = %+v\n", resp)
	fmt.Printf("err = %+v\n", err)
	fmt.Println(buildHeaderString(resp.Header))
}

func TestHttpPost() {
	p := &Person{"hello", "world", time.Now()}
	pData, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	res, err := http.Post("http://127.0.0.1:8080/testHttpPost", "application/json", bytes.NewBuffer(pData))
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	fmt.Printf("string(data) = %+v\n", string(data))
}

func buildHeaderString(header http.Header) string {
	lines := make([]string, 0)
	for k, v := range header {
		lines = append(lines, fmt.Sprintf("%s: %s", k, strings.Join(v, "")))
	}
	return strings.Join(lines, "\n")
}
