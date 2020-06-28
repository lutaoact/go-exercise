package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
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
	prefixesKey           = "archive:prefixes" // set: prefix
	prefixesKeyBak        = "archive:prefixes:bak"
	markerKey             = "archive:marker"        // hash: prefix => marker
	statResultKey         = "archive:statResult:%s" // hash
	finalResultKey        = "archive:finalResult"   // hash
	compareCountKey       = "archive:compareCount"  // integer
	count           int64 = 1024
)

var (
	bucket         *oss.Bucket
	redisClient    *redis.Client
	awsBucket      *string
	svc            *s3.S3
	errLogger      = logrus.New()
	diffLogger     = logrus.New()
	wg             = &sync.WaitGroup{}
	diffWg         = &sync.WaitGroup{}
	limitChan      = make(chan struct{}, 512)
	groupNameMap   = make(map[int]string)
	compareKeyChan = make(chan string, 100)
)

func init() {
	_ = flag.Set("logtostderr", "true")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	initRedis()
	initOSS()
	initAWS()
	setLogger()
	initGroupNameMap()
}

func main() {
	redisClient = redisdb.GetRedis()
	bucket = myoss.Bucket()

	go func() {
		for key := range compareKeyChan {
			compareOneKeyBetweenTwoCloud(key)
		}
	}()

	prefixChan := scanPrefixes()
	for prefix := range prefixChan {
		fmt.Printf("prefix = %+v\n", prefix)
		comeIn()
		go processOnePrefix(prefix)
	}

	wg.Wait()
	diffWg.Wait()
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

func addKeyForCompare(key string) {
	if rand.Float32() < 0.5 {
		return
	}

	select {
	case compareKeyChan <- key:
		diffWg.Add(1)
	default:
		fmt.Printf("ignore key = %+v\n", key)
	}
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

		go addKeyForCompare(randomGetOneKey(lor.Objects))

		statResult := make(map[int]map[string]int64)
		for _, obj := range lor.Objects {
			tuple := computeSizeTuple(obj.Size)
			if _, ok := statResult[tuple.Group]; ok {
				statResult[tuple.Group]["block_count"] += tuple.BlockCount
				statResult[tuple.Group]["file_count"] += 1
				statResult[tuple.Group]["size"] += obj.Size
			} else {
				statResult[tuple.Group] = make(map[string]int64)
				statResult[tuple.Group]["block_count"] = tuple.BlockCount
				statResult[tuple.Group]["file_count"] = 1
				statResult[tuple.Group]["size"] = obj.Size
			}
		}

		if err := addStatResult(statResult); err != nil {
			log.Fatalf("addStatResult: %+v %+v %+v", err, marker, prefix)
		}

		// TODO 以上统计结果和lor.NextMarker存储到redis
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

func addStatResult(statResult map[int]map[string]int64) error {
	pipe := redisClient.Pipeline()
	for groupNum, stat := range statResult {
		groupName := groupNameMap[groupNum]
		key := fmt.Sprintf(statResultKey, groupName)

		for hashkey, value := range stat {
			pipe.HIncrBy(key, hashkey, value)
			pipe.HIncrBy(finalResultKey, hashkey, value)
		}
	}
	_, err := pipe.Exec()

	return err
}

func clearPrefix(prefix string) error {
	pipe := redisClient.Pipeline()
	pipe.SRem(prefixesKey, prefix)
	pipe.SAdd(prefixesKeyBak, prefix)
	pipe.HDel(markerKey, prefix)
	_, err := pipe.Exec()
	return err
}

type sizeTuple struct {
	Size       int64
	BlockCount int64
	Group      int
}

func computeSizeTuple(size int64) *sizeTuple {
	n := size >> 10 // KB为单位
	i := 0
	for n > 0 && i < 10 {
		i += 1
		n = n >> 1
	}

	ret := &sizeTuple{size, 0, 1 << (10 + i)}
	if i <= 6 {
		ret.BlockCount = 1
	} else {
		ret.BlockCount = (size >> 16) + 1
	}

	return ret
}

func randomGetOneKey(objs []oss.ObjectProperties) string {
	index := rand.Int63n(int64(len(objs)))
	return objs[index].Key
}

func compareOneKeyBetweenTwoCloud(key string) {
	defer diffWg.Done()

	ossChan := asyncCall(getObjFromOSS, key)
	awsChan := asyncCall(getObjFromAWS, key)
	ossRet := <-ossChan
	awsRet := <-awsChan
	if ossRet.Err != nil || awsRet.Err != nil {
		errLogger.WithFields(logrus.Fields{
			"key":     key,
			"oss_err": fmt.Sprintf("%+v", ossRet.Err),
			"aws_err": fmt.Sprintf("%+v", awsRet.Err),
		}).Errorln()
		return
	}

	if bytes.Compare(ossRet.Ret.([]byte), awsRet.Ret.([]byte)) != 0 {
		diffLogger.WithFields(logrus.Fields{
			"key": key,
		}).Infoln()
		return
	}
	fmt.Println("same content:", key)
	// 计个数，错误不重要，所以直接忽略了
	if err := redisClient.Incr(compareCountKey).Err(); err != nil {
		log.Fatal("incr compareCountKey: %+v", err)
	}
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
	errFile, err := os.OpenFile("useraudio.err.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	redisdb.MustInitRedisClient(&common.RedisConfig{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   31,
		//Addr:     "127.0.0.1:6379",
		//DB:       3,
		Password: "",
	})
}

func initOSS() {
	myoss.MustInit(&myoss.OSSConfig{
		Bucket:   "archive-useraudio",
		Endpoint: "oss-cn-shanghai-internal.aliyuncs.com",
		//Endpoint:  "oss-cn-shanghai.aliyuncs.com",
		KeyID:     os.Getenv("USERAUDIO_AK"),
		KeySecret: os.Getenv("USERAUDIO_SK"),
	})
}

func initAWS() {
	awsBucket = aws.String("useraudio")
	sess := session.Must(session.NewSession())
	svc = s3.New(sess)
}

func initGroupNameMap() {
	groupNameMap[1<<10] = "1KB"
	groupNameMap[1<<11] = "2KB"
	groupNameMap[1<<12] = "4KB"
	groupNameMap[1<<13] = "8KB"
	groupNameMap[1<<14] = "16KB"
	groupNameMap[1<<15] = "32KB"
	groupNameMap[1<<16] = "64KB"
	groupNameMap[1<<17] = "128KB"
	groupNameMap[1<<18] = "256KB"
	groupNameMap[1<<19] = "512KB"
	groupNameMap[1<<20] = "1MB"
}
