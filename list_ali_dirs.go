package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis"

	"git.llsapp.com/algapi/alix/pkg/common"
	"git.llsapp.com/algapi/alix/pkg/redisdb"
	"git.llsapp.com/algapi/sonic/pkg/myoss"
)

var (
	bucket      *oss.Bucket
	redisClient *redis.Client
)

func init() {
	_ = flag.Set("logtostderr", "true")
	flag.Parse()

	redisdb.MustInitRedisClient(&common.RedisConfig{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   31,
		//Addr:     "127.0.0.1:6379",
		//DB:       3,
		Password: "",
	})

	myoss.MustInit(&myoss.OSSConfig{
		Bucket:   "archive-useraudio",
		Endpoint: "oss-cn-shanghai-internal.aliyuncs.com",
		//Endpoint:  "oss-cn-shanghai.aliyuncs.com",
		KeyID:     oss.Getenv("USERAUDIO_AK"),
		KeySecret: oss.Getenv("USERAUDIO_SK"),
	})
}

func main() {
	redisClient = redisdb.GetRedis()
	bucket = myoss.Bucket()

	outChan := ListDirs()
	for prefixes := range outChan {
		fmt.Printf("len(prefixes) = %+v\n", len(prefixes))
		err := redisClient.SAdd("archive:prefixes", convert(prefixes)...).Err()
		if err != nil {
			log.Printf("SAdd: %+v %+v", err, prefixes)
		}
	}
}

func convert(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func ListDirs() <-chan []string {
	outChan := make(chan []string, 10)

	go func() {
		delimiter := oss.Delimiter("/")
		marker := oss.Marker("")
		for {
			// MaxKeys最大值为1000，但若设为1000，会导致超时
			lor, err := bucket.ListObjects(oss.MaxKeys(500), marker, delimiter)
			if err != nil {
				log.Fatal("ListByPrefix: %+v %+v", err, marker)
			}

			marker = oss.Marker(lor.NextMarker)
			outChan <- lor.CommonPrefixes
			fmt.Printf("lor = %+v\n", lor)

			if !lor.IsTruncated {
				close(outChan)
				break
			}
		}
	}()

	return outChan
}
