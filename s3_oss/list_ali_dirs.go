package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis"

	"git.llsapp.com/algapi/alix/pkg/common"
	"git.llsapp.com/algapi/alix/pkg/redisdb"
	"git.llsapp.com/algapi/sonic/pkg/myoss"
)

var (
	bucket        *oss.Bucket
	redisClient   *redis.Client
	darwinDirsKey = "archive:darwin:dirs"
)

func init() {
	_ = flag.Set("logtostderr", "true")
	flag.Parse()

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalln("Failed to parse REDIS_DB")
	}

	redisdb.MustInitRedisClient(&common.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		DB:       redisDB,
		Password: "",
	})

	myoss.MustInit(&myoss.OSSConfig{
		Bucket:    "archive-useraudio",
		Endpoint:  os.Getenv("ENDPOINT"),
		KeyID:     os.Getenv("USERAUDIO_AK"),
		KeySecret: os.Getenv("USERAUDIO_SK"),
	})
}

func main() {
	redisClient = redisdb.GetRedis()
	bucket = myoss.Bucket()

	outChan := ListDirs("darwin/")
	for prefixes := range outChan {
		fmt.Printf("len(prefixes) = %+v\n", len(prefixes))
		err := redisClient.SAdd(darwinDirsKey, convert(prefixes)...).Err()
		if err != nil {
			log.Fatalf("SAdd: %+v %+v", err, prefixes)
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

func ListDirs(prefix string) <-chan []string {
	outChan := make(chan []string, 10)

	go func() {
		delimiter := oss.Delimiter("/")
		marker := oss.Marker("")
		for {
			// MaxKeys最大值为1000，但若设为1000，会导致超时
			lor, err := bucket.ListObjects(oss.Prefix(prefix), oss.MaxKeys(500), marker, delimiter)
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
