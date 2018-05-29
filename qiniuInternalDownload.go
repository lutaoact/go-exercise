package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lutaoact/go-exercise/secret"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/sirupsen/logrus"
)

var accessKey = secret.AK
var secretKey = secret.SK
var mac = qbox.NewMac(accessKey, secretKey)

var path = "ke/docker/registry/v2/blobs/sha256/94/941f946e34a099e4323c0a2180c2b5cc38049dfe0917f9e4db9fe5b2814a1510/data"

func main() {
	err := download(privateURL(path))
	fmt.Println(err)
}

func privateURL(path string) string {
	domain := "http://p7b7qb2jj.bkt.clouddn.com"
	key := path
	fmt.Printf("key = %+v\n", key)
	deadline := time.Now().Add(time.Second * 3600 * 3).Unix() //3小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)

	return privateAccessURL
}

// curl -v -H 'Host: p7b7qb2jj.bkt.clouddn.com' 'http://xsio.qiniu.io/ke/docker/registry/v2/repositories/lutaoact/bigimage/_uploads/06348fc6-82b8-4b8c-83a6-9f03693f5a45/startedat?e=1527484460&token=zSPmaQUCLkSRWtjNBTJdHFi59cxr6iJBypL-tVoK:7ePvjFnpoOmX7xY881Z-S8hkhDo='
func download(privateAccessURL string) error {
	//os.Setenv("HTTP_PROXY", "http://localhost:8888")

	logger := logrus.WithFields(logrus.Fields{
		"func": "download", "privateAccessURL": privateAccessURL,
	})

	fmt.Printf("privateAccessURL = %+v\n", privateAccessURL)
	privateAccessURL = strings.Replace(
		privateAccessURL, "p7b7qb2jj.bkt.clouddn.com", "xsio.qiniu.io", 1,
	)

	fmt.Printf("privateAccessURL = %+v\n", privateAccessURL)
	//req, err := http.NewRequest("GET", privateAccessURL, nil)
	//if err != nil {
	//	logger.Error("http.NewRequest:", err)
	//	return err
	//}

	u, _ := url.Parse(privateAccessURL)
	req := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       nil,
		Host:       u.Host,
	}

	//fmt.Printf("req.Header = %+v\n", req.Header)

	//t := NewHostTransport("p7b7qb2jj.bkt.clouddn.com", nil)
	//client := &http.Client{Transport: t}
	//resp, err := client.Do(req)
	req.Host = "p7b7qb2jj.bkt.clouddn.com"
	//req.URL.Host = "xsio.qiniu.io"
	//req.Host = "xsio.qiniu.io"

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("http.Do:", err)
		return err
	}
	defer resp.Body.Close()

	fp, err := os.OpenFile("/tmp/qiniuInternalDownload", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("OpenFile:", "/tmp/qiniuInternalDownload")
		return err
	}

	writer := bufio.NewWriter(fp)
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		logger.Error("io.Copy:", "/tmp/qiniuInternalDownload")
		return err
	}

	if err := writer.Flush(); err != nil {
		logger.Error(err)
		return err
	}

	if err := fp.Sync(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

type HostTransport struct {
	Host      string
	Transport http.RoundTripper
}

func (t *HostTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("Host", t.Host)
	return t.Transport.RoundTrip(req)
}

func NewHostTransport(host string, transport http.RoundTripper) *HostTransport {
	//	if transport == nil {
	//		transport = http.DefaultTransport
	//	}

	return &HostTransport{Transport: http.DefaultTransport, Host: host}
}
