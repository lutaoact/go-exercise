package main

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	//log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//		log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)
}

func main() {
	ak, sk := "ak:xxxx", "sk:yyyy"
	logger := log.WithFields(log.Fields{"ak": ak, "sk": sk})
	err := errors.New("this is a error")
	logger.Error("getb error:", err)

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Debug("A group of walrus emerges from the ocean")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Error("The ice breaks!")

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})
	contextLogger.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Error("The ice breaks!")

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")
}
