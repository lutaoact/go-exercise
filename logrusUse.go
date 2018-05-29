package main

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetFormatter(&logrus.TextFormatter{})

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
		FullTimestamp:   true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only logrus the warning severity or above.
	//		logrus.SetLevel(logrus.WarnLevel)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	ak, sk := "ak:xxxx", "sk:yyyy"
	logger := logrus.WithFields(logrus.Fields{"ak": ak, "sk": sk})
	err := errors.New("this is a error")
	logger.Error("getb error:", err)

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Debug("A group of walrus emerges from the ocean")

	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 100,
	}).Error("The ice breaks!")

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := logrus.WithFields(logrus.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})
	contextLogger.WithFields(logrus.Fields{
		"omg":    true,
		"number": 100,
	}).Error("The ice breaks!")

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too", "hh")

	logrus.Errorf("%+v", &MyPutRet{Key: "hhhh"})
}
