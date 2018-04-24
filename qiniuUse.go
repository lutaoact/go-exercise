package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/lutaoact/go-exercise/secret"
	"github.com/qbox/stark/hub/redisdb"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
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
var cfg = storage.Config{UseHTTPS: false}

// 指定空间所在的区域，如果不指定将自动探测
// 如果没有特殊需求，默认不需要指定
//cfg.Zone=&storage.ZoneHuabei
var bucketManager = storage.NewBucketManager(mac, &cfg)

func main() {
	//	formUpload()
	//resumeUpload()
	//stat()
	list()
}

func list() {
	prefix := ""

	//这个字段的含义比较模糊，我这里解释一下
	//ListFiles在默认情况下是不返回目录的，但是如果指定了delimiter，
	//就会把从prefix开始，到delimiter为止的部分，看做是目录，
	//一般情况下delimiter就是/，但其实可以随便指定，我指定一个字符m也是可以的
	delimiter := "m"
	//初始列举marker为空
	marker := ""
	limit := 10
	entries, commonPrefixes, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, prefix, delimiter, marker, limit)
	fmt.Println(entries, commonPrefixes, nextMarker, hasNext, err)
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
	config := redisdb.Config{Addr: "127.0.0.1:6379"}
	redisClient := redisdb.InitRedis(&config)
	fmt.Println(redisClient.Ping())

	keyToOverwrite := "random2.data"
	localFile := "random2.data"
	key := "random2.data"

	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", bucket, keyToOverwrite), //覆盖上传
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
