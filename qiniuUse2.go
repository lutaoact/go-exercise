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
var cfg = storage.Config{UseHTTPS: false}

// 指定空间所在的区域，如果不指定将自动探测
// 如果没有特殊需求，默认不需要指定
//cfg.Zone=&storage.ZoneHuabei
var bucketManager = storage.NewBucketManager(mac, &cfg)

func main() {
	resumeUpload2()
	resumeUpload()
}

func resumeUpload() {

	var accessKey = secret.AK
	var secretKey = secret.SK
	var mac = qbox.NewMac(accessKey, secretKey)
	var bucket = "registry-lutao"

	//keyToOverwrite := "ke/docker/registry/v2/repositories/lutaoact/hello-world/_uploads/5083e8ec-5eea-448b-97af-77546db9ec79/startedat"
	localFile := "/tmp/docker/registry/v2/repositories/lutaoact/hello-world/_uploads/5083e8ec-5eea-448b-97af-77546db9ec79/startedat"
	key := "ke/docker/registry/v2/repositories/lutaoact/hello-world/_uploads/5083e8ec-5eea-448b-97af-77546db9ec79/startedat"

	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", bucket, key), //覆盖上传
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
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
