package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const maxKeys = 10

var pageNum = 0

func processObjects(output *s3.ListObjectsV2Output, lastPage bool) bool {
	pageNum++
	//fmt.Printf("output = %+v\n", output)
	fmt.Printf("output.NextContinuationToken = %+v\n", *output.NextContinuationToken)
	fmt.Printf("lastPage = %+v\n", lastPage)
	if lastPage {
		return false
	}

	fmt.Printf("*output.KeyCount = %+v\n", *output.KeyCount)
	index := rand.Int63n(*output.KeyCount)
	fmt.Printf("index = %+v\n", index)
	obj := output.Contents[index]
	fmt.Printf("obj = %+v\n", obj)
	// 正确的return
	// return !lastPage

	return pageNum <= 3
}

func main() {
	rand.Seed(time.Now().UnixNano())
	awsBucket := aws.String("useraudio")
	//aliBucket := "archive-useraudio"

	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	input := &s3.ListObjectsV2Input{
		Bucket:  awsBucket,
		MaxKeys: aws.Int64(maxKeys),
	}

	err := svc.ListObjectsV2Pages(input, processObjects)
	fmt.Println(err)
}

type chanRet struct {
	Ret interface{}
	Err error
}

func diffOneKey(wg *sync.WaitGroup, obj *s3.Object, c chan<- *retForChan) {
	defer wg.Done()

}
