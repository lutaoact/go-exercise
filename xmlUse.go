package main

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/aws/aws-sdk-go/service/s3"
)

var data []byte

func main() {
	initData()

	var obj s3.DeleteObjectsInput
	decoder := xml.NewDecoder(bytes.NewReader(data))
	err := xmlutil.UnmarshalXML(&obj, decoder, "")
	fmt.Println(obj, err)
}

func initData() {
	data = []byte(``)
}
