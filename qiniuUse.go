package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/lutaoact/go-exercise/secret"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/sirupsen/logrus"
)

type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func md5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type ProgressRecord struct {
	Progresses []storage.BlkputRet `json:"progresses"`
}

var accessKey = secret.AK
var secretKey = secret.SK
var mac = qbox.NewMac(accessKey, secretKey)

var bucket = "registry-lutao"

//var bucket = "registry-qiniu-com"
var cfg = storage.Config{UseHTTPS: false}

// 指定空间所在的区域，如果不指定将自动探测
// 如果没有特殊需求，默认不需要指定
//cfg.Zone=&storage.ZoneHuabei
var bucketManager = storage.NewBucketManager(mac, &cfg)

func main() {
	//	formUpload()
	//resumeUpload()
	//resumeUpload2()
	//stat()
	//list()
	get()
	//delete()
}

func delete() {
	//key := "random.data"
	key := "ke/docker/registry/v2/blobs/sha256/0b/0b1edfbffd27c935a666e233a0042ed634205f6f754dbe20769a60369c614f85/data"
	err := bucketManager.Delete(bucket, key)
	if err == nil {
		return
	}
	fmt.Printf("err = %+v\n", err)

	err1 := err.(*storage.ErrorInfo)
	fmt.Printf("err = %+v\n", err1)
	fmt.Printf("err = %+v\n", err1.Err)
	fmt.Printf("err = %+v\n", err1.Key)
	fmt.Printf("err = %+v\n", err1.Reqid)
	fmt.Printf("err = %+v\n", err1.Errno)
	fmt.Printf("err = %+v\n", err1.Code)
}

func get() {
	//domain := "http://p7b7qb2jj.bkt.clouddn.com"
	domain := "http://reg-store.kirkcloud.com"
	key := "ke/docker/registry/v2/blobs/sha256/00/00008c158f626dbc5b00ee179a4e2122cb8a58433a161d41affdc9dc2e9dc418/data"
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)

	//	client := &http.Client{}

	req, err := http.NewRequest("GET", privateAccessURL, nil)
	if err != nil {
		logrus.Error(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		logrus.Error(err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}

	err = ioutil.WriteFile("/tmp/file.data", data, 0644)
	fmt.Println(err)

	resp.Body.Close()
}

func list() {
	//prefix := "ke/docker/registry/v2/repositories/kirk-apps/redis-controller/_uploads/e68aed42-9d36-4ac9-b8e3-d563d2d9546"
	//prefix := "ke/docker/registry/v2/repositories/kirk-apps/redis-controller/_uploads/e68aed42-9d36-4ac9-b8e3-d563d2d9546a"
	prefix := "ke/docker/registry/v2/repositories/spock-release-candidates/rtn-frontend/_uploads/b7fc2462-1786-4081-9ca0-a01e227951fc"

	//这个字段的含义比较模糊，我这里解释一下
	//ListFiles在默认情况下是不返回目录的，但是如果指定了delimiter，
	//就会把从prefix开始，到delimiter为止的部分，看做是目录，
	//一般情况下delimiter就是/，但其实可以随便指定，我指定一个字符m也是可以的
	delimiter := ""
	//初始列举marker为空
	marker := ""
	limit := 10
	entries, commonPrefixes, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, prefix, delimiter, marker, limit)
	fmt.Println(entries, commonPrefixes, nextMarker, hasNext, err)
	for _, entry := range entries {
		fmt.Printf("entry.Key = %+v\n", entry.Key)
	}
}

//不支持Stat一个目录，只能stat文件
func stat() {
	cfg := storage.Config{
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	//fileInfo, err := bucketManager.Stat(bucket, "random.data")
	fileInfo, err := bucketManager.Stat(bucket, "health/out.txt")
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("fileInfo = %+v\n", fileInfo)
	fmt.Println(fileInfo.String())
	//可以解析文件的PutTime
	fmt.Println(storage.ParsePutTime(fileInfo.PutTime))
}

func resumeUpload() {
	//	config := redisdb.Config{Addr: "127.0.0.1:6379"}
	//	redisClient := redisdb.InitRedis(&config)
	//	fmt.Println(redisClient.Ping())

	//	keyToOverwrite := "ke/docker/registry/v2/repositories/lutaoact/hello-world/_uploads/5083e8ec-5eea-448b-97af-77546db9ec79/startedat"
	localFile := "/tmp/docker/registry/v2/repositories/lutaoact/hello-world/_uploads/5083e8ec-5eea-448b-97af-77546db9ec79/startedat"
	key := "ke/docker/registry/v2/repositories/lutaoact/hello-world/_uploads/5083e8ec-5eea-448b-97af-77546db9ec79/startedat"

	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", bucket, key), //覆盖上传
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	resumeUploader := storage.NewResumeUploader(&cfg)
	ret := MyPutRet{}

	fileInfo, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Println(statErr)
		return
	}
	fmt.Printf("fileInfo = %+v\n", fileInfo)
	fileSize := fileInfo.Size()

	fmt.Printf("BlockCount = %+v\n", storage.BlockCount(fileSize))

	progressRecord := ProgressRecord{}
	progressRecord.Progresses = make(
		[]storage.BlkputRet,
		storage.BlockCount(fileSize),
	)

	progressLock := sync.RWMutex{}

	putExtra := storage.RputExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
		ChunkSize: 2 * 1024 * 1024, //分片的大小
		Notify: func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			progressLock.Lock()
			defer progressLock.Unlock()
			//将进度序列化，然后写入文件
			progressRecord.Progresses[blkIdx] = *ret
			progressBytes, _ := json.Marshal(progressRecord)
			fmt.Printf("string(progressBytes) = %+v\n", string(progressBytes))
			fmt.Println("upload progress:", blkIdx, blkSize)
		},
	}
	err := resumeUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Bucket, ret.Key, ret.Fsize, ret.Hash, ret.Name)

	//	fmt.Println("resumeUpload2===============")
	//	fmt.Printf("resumeUploader = %+v\n", resumeUploader)
	//	resumeUpload2(&cfg)
}

func resumeUpload2() {
	var accessKey = secret.AK
	var secretKey = secret.SK
	var mac = qbox.NewMac(accessKey, secretKey)
	var bucket = "registry-lutao"

	//keyToOverwrite := "random3.data"
	localFile := "./random3.data"
	key := "random3.data"

	cfg := &storage.Config{UseHTTPS: false}

	ret := MyPutRet{}
	//	putExtra := storage.RputExtra{}

	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", bucket, key), //覆盖上传
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	upToken := putPolicy.UploadToken(mac)

	uploader := storage.NewResumeUploader(cfg)

	fileInfo, statErr := os.Stat(localFile)
	if statErr != nil {
		fmt.Println(statErr)
		return
	}
	fmt.Printf("fileInfo = %+v\n", fileInfo)
	fileSize := fileInfo.Size()

	fmt.Printf("BlockCount = %+v\n", storage.BlockCount(fileSize))

	progressRecord := ProgressRecord{}
	progressRecord.Progresses = make(
		[]storage.BlkputRet,
		storage.BlockCount(fileSize),
	)

	progressLock := sync.RWMutex{}

	putExtra := storage.RputExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
		Notify: func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			progressLock.Lock()
			defer progressLock.Unlock()
			//将进度序列化，然后写入文件
			progressRecord.Progresses[blkIdx] = *ret
			progressBytes, _ := json.Marshal(progressRecord)
			fmt.Printf("string(progressBytes) = %+v\n", string(progressBytes))
			fmt.Println("upload progress:", blkIdx, blkSize)
		},
	}

	err := uploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)

	if err != nil {
		logrus.Error("2 PutFile:", err)
		return
	}
	fmt.Println(ret.Bucket, ret.Key, ret.Fsize, ret.Hash, ret.Name)
}

func formUpload() {
	keyToOverwrite := "random.data"
	putPolicy := storage.PutPolicy{
		//		Scope:      bucket,
		Scope:      fmt.Sprintf("%s:%s", bucket, keyToOverwrite), //覆盖上传
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	//	putPolicy := storage.PutPolicy{Scope: fmt.Sprintf("%s:%s", bucket, keyToOverwrite)}

	upToken := putPolicy.UploadToken(mac)
	fmt.Printf("upToken = %+v\n", upToken)

	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := MyPutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	localFile := "random.data"
	key := "random.data"
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Bucket, ret.Key, ret.Fsize, ret.Hash, ret.Name)
}
