package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	awsBucket *string
	svc       *s3.S3
	keys      []string
)

func init() {
	initAWS()
	initKeys()
}

func main() {
	ids := make([]*s3.ObjectIdentifier, len(keys))
	for i, key := range keys {
		ids[i] = &s3.ObjectIdentifier{
			Key: aws.String(key),
		}
	}

	input := &s3.DeleteObjectsInput{
		Bucket: awsBucket,
		Delete: &s3.Delete{
			Objects: ids,
			Quiet:   aws.Bool(true),
		},
	}
	//fmt.Printf("input = %+v\n", input)
	output, err := svc.DeleteObjects(input)
	fmt.Println(output, err)
}

func initAWS() {
	awsBucket = aws.String("useraudio")
	sess := session.Must(session.NewSession())
	svc = s3.New(sess, &aws.Config{LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody)})
}

func initKeys() {
	keys = []string{
		"2.pb",
	}
}
