package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	"git.llsapp.com/algapi/alix/pkg/common"
	"git.llsapp.com/algapi/alix/pkg/redisdb"
	"git.llsapp.com/algapi/sonic/pkg/myoss"
)

const (
	markerKey             = "archive:deleting:marker" // hash: prefix => marker
	statResultKey         = "archive:statResult:%s"   // hash
	finalResultKey        = "archive:finalResult"     // hash
	compareCountKey       = "archive:compareCount"    // integer
	count           int64 = 1024
)

var (
	bucket      *oss.Bucket
	redisClient *redis.Client
	awsBucket   *string
	svc         *s3.S3
	errLogger   = logrus.New()
	diffLogger  = logrus.New()
	wg          = &sync.WaitGroup{}
	diffWg      = &sync.WaitGroup{}
	limitChan   = make(chan struct{}, 1024)

	// prefixesKey通过参数传递，可以是 archive:prefixes:deleting 或 archive:darwin:dirs:deleting
	// prefixesKeyBak为 prefixesKey + ":bak"
	prefixesKey    string // set: prefix
	prefixesKeyBak string // zset: prefix => ts
)

func init() {
	_ = flag.Set("logtostderr", "true")
	flag.Parse()

	if len(os.Args) != 2 || (os.Args[1] != "archive:prefixes:deleting" && os.Args[1] != "archive:darwin:dirs:deleting") {
		parts := strings.Split(os.Args[0], "/")
		fmt.Printf("Usage: ./%s [prefixes]\n", parts[len(parts)-1])
		fmt.Println(`prefixes can be: "archive:prefixes:deleting" or "archive:darwin:dirs:deleting"`)
		os.Exit(1)
	}

	prefixesKey = os.Args[1]
	prefixesKeyBak = prefixesKey + ":bak"

	rand.Seed(time.Now().UnixNano())

	initRedis()
	initOSS()
	initAWS()
	setLogger()
}

func main() {
	redisClient = redisdb.GetRedis()
	bucket = myoss.Bucket()

	prefixChan := scanPrefixes()
	for prefix := range prefixChan {
		fmt.Printf("prefix = %+v\n", prefix)
		comeIn()
		go processOnePrefix(prefix)
	}

	wg.Wait()
	diffWg.Wait()
	fmt.Println("done")
}

func comeIn() {
	wg.Add(1)
	limitChan <- struct{}{}
}

func getOut() {
	wg.Done()
	<-limitChan
}

func scanPrefixes() <-chan string {
	outChan := make(chan string, count)

	go func() {
		var cursor uint64 = 0
		for {
			prefixes, nextCursor, err := redisClient.SScan(prefixesKey, cursor, "", count).Result()
			if err != nil {
				log.Fatal("sscan: %+v %+v", err, cursor)
			}

			fmt.Println(prefixes, nextCursor)

			for _, v := range prefixes {
				outChan <- v
			}
			if nextCursor == 0 {
				close(outChan)
				break
			}

			cursor = nextCursor
		}
	}()
	return outChan
}

func processOnePrefix(prefix string) {
	defer getOut()

	// 根据prefix从redis中获取marker，这个marker就是该prefix的处理进度
	marker, err := redisClient.HGet(markerKey, prefix).Result()
	if err != nil && err != redis.Nil {
		log.Fatalf("hget: %+v %+v %+v", err, marker, prefix)
	}

	lorChan := myoss.ListByPrefixAndMarker(prefix, marker)
	for lor := range lorChan {
		fmt.Printf("len(objs) = %+v\n", len(lor.Objects))

		// 批量删除文件
		batchDeleteObjects(lor.Objects)

		// lor.NextMarker存储到redis，用来标记下次启动时的状态
		err := redisClient.HSet(markerKey, prefix, lor.NextMarker).Err()
		if err != nil {
			log.Fatalf("hset markerKey: %+v %+v %+v", err, marker, prefix)
		}
	}

	// 删除redis集合prefixesKey中的该prefix，并存储到备份集合中
	// 删除markerKey中的该prefix
	if err = clearPrefix(prefix); err != nil {
		log.Fatalf("pipe.Exec(): %+v %+v", err, prefix)
	}
}

func batchDeleteObjects(objs []oss.ObjectProperties) error {
	ids := make([]*s3.ObjectIdentifier, len(objs))
	for i, obj := range objs {
		ids[i] = &s3.ObjectIdentifier{
			Key: aws.String(obj.Key),
		}
	}

	input := &s3.DeleteObjectsInput{
		Bucket: awsBucket,
		Delete: &s3.Delete{
			Objects: ids,
			Quiet:   aws.Bool(true),
		},
	}
	// 忽略错误，先删除那些可以删除的
	output, err := svc.DeleteObjects(input)
	if err != nil {
		errLogger.WithFields(logrus.Fields{"input": input, "err": err.Error()}).Errorln()
		return nil
	}
	if len(output.Errors) > 0 {
		errLogger.WithFields(logrus.Fields{"input": input, "output": output}).Errorln()
		return nil
	}
	log.Printf("delete success len(ids): %+v", len(ids))
	return nil
}

func clearPrefix(prefix string) error {
	pipe := redisClient.Pipeline()
	pipe.SRem(prefixesKey, prefix)
	// 记录删除完成时间
	pipe.ZAdd(prefixesKeyBak, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: prefix,
	})
	pipe.HDel(markerKey, prefix)
	_, err := pipe.Exec()
	return err
}

type retForChan struct {
	Ret interface{}
	Err error
}

type MyFunc func(args ...interface{}) (interface{}, error)

func asyncCall(fn MyFunc, params ...interface{}) <-chan *retForChan {
	out := make(chan *retForChan)
	go func() {
		ret := new(retForChan)
		ret.Ret, ret.Err = fn(params...)
		out <- ret
	}()
	return out
}

func getObjFromOSS(args ...interface{}) (interface{}, error) {
	key := args[0].(string)

	body, err := bucket.GetObject(key)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return ioutil.ReadAll(body)
}

func getObjFromAWS(args ...interface{}) (interface{}, error) {
	key := args[0].(string)

	input := &s3.GetObjectInput{Bucket: awsBucket, Key: aws.String(key)}
	obj, err := svc.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()
	return ioutil.ReadAll(obj.Body)
}

func setLogger() {
	errFile, err := os.OpenFile("deleting.err.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to log to file, using default stderr")
	}
	errLogger.Out = errFile

	diffFile, err := os.OpenFile("useraudio.diff.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to log to file, using default stderr")
	}
	diffLogger.Out = diffFile
}

func initRedis() {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalln("Failed to parse REDIS_DB")
	}

	redisdb.MustInitRedisClient(&common.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		DB:       redisDB,
		Password: "",
	})
}

func initOSS() {
	myoss.MustInit(&myoss.OSSConfig{
		Bucket:    "archive-useraudio",
		Endpoint:  os.Getenv("ENDPOINT"),
		KeyID:     os.Getenv("USERAUDIO_AK"),
		KeySecret: os.Getenv("USERAUDIO_SK"),
	})
}

func initAWS() {
	awsBucket = aws.String("useraudio")
	sess := session.Must(session.NewSession())
	svc = s3.New(sess, &aws.Config{LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody)})
}
