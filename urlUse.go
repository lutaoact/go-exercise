package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lutaoact/go-exercise/secret"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var accessKey = secret.AK
var secretKey = secret.SK
var mac = qbox.NewMac(accessKey, secretKey)

func main() {
	u, err := url.Parse("http://p7b7qb2jj.bkt.clouddn.com")
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("u = %+v\n", u)
	fmt.Printf("u.Host = %+v\n", u.Host)
	fmt.Printf("url = %+v\n", privateURL("/docker/registry/v2/blobs/sha256/00/0074283d71418f7409e1905b287d8c900db7415d064b46cb714f1511c8a1078b/data"))
}

func privateURL(path string) string {
	domain := "http://p7b7qb2jj.bkt.clouddn.com"
	key := kodoKey(path)
	fmt.Printf("key = %+v\n", key)
	deadline := time.Now().Add(time.Second * 3600 * 3).Unix() //3小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)

	return privateAccessURL
}

func kodoKey(path string) string {
	return strings.TrimLeft(strings.TrimRight("ke", "/")+path, "/")
}
