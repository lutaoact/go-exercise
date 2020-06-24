package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
	//domain := "http://xsio.qiniu.io"
	domain := "http://p7b7qb2jj.bkt.clouddn.com"
	key := path
	fmt.Printf("key = %+v\n", key)
	deadline := time.Now().Add(time.Second * 3600 * 3).Unix() //3小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)

	return privateAccessURL
}

// curl -v -H 'Host: p7b7qb2jj.bkt.clouddn.com' 'http://xsio.qiniu.io/ke/docker/registry/v2/repositories/lutaoact/bigimage/_uploads/06348fc6-82b8-4b8c-83a6-9f03693f5a45/startedat?e=1527484460&token=zSPmaQUCLkSRWtjNBTJdHFi59cxr6iJBypL-tVoK:7ePvjFnpoOmX7xY881Z-S8hkhDo='
// curl -v -H 'Host: p7b7qb2jj.bkt.clouddn.com' 'http://xsio.qiniu.io/ke/docker/registry/v2/blobs/sha256/94/941f946e34a099e4323c0a2180c2b5cc38049dfe0917f9e4db9fe5b2814a1510/data?e=1527498010&token=zSPmaQUCLkSRWtjNBTJdHFi59cxr6iJBypL-tVoK:huXx9mLGY_NIjGaZ-3d3cBiSDgs='
func download(privateAccessURL string) error {
	//os.Setenv("HTTP_PROXY", "http://localhost:8888")
	logger := logrus.WithFields(logrus.Fields{
		"func": "download", "privateAccessURL": privateAccessURL,
	})
	privateAccessURL = strings.Replace(
		privateAccessURL, "p7b7qb2jj.bkt.clouddn.com", "xsio.qiniu.io", 1,
	)
	fmt.Printf("privateAccessURL = %+v\n", privateAccessURL)

	req, err := http.NewRequest("GET", privateAccessURL, nil)
	if err != nil {
		logger.Error("http.NewRequest:", err)
		return err
	}

	req.Host = "p7b7qb2jj.bkt.clouddn.com"
	//req.Host = "xsio.qiniu.io"

	res, err := http.DefaultClient.Do(req)
	//res, err := rpc.DefaultClient.DoRequest(nil, "GET", "-H p7b7qb2jj.bkt.clouddn.com "+privateAccessURL)
	if err != nil {
		logger.Error("http.Do:", err)
		return err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("ReadAll:", err)
		return err
	}
	fmt.Printf("string = %+v\n", string(data))
	return nil
}
