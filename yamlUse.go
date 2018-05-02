package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	file, err := os.Open("./kodo-dev.yml")
	if err != nil {
		logrus.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("string(data) = %+v\n", string(data))
	result := map[string]interface{}{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(err)
	fmt.Println(result["storage"])

	storage := result["storage"].(map[interface{}]interface{})
	mymap := map[string]interface{}{}
	for k, v := range storage {
		mymap[k.(string)] = v
	}
	fmt.Println(mymap["kodo"])
	kodo := mymap["kodo"].(map[interface{}]interface{})
	kodoMap := map[string]interface{}{}
	for k, v := range kodo {
		kodoMap[k.(string)] = v
	}
	fmt.Println(kodoMap["rootdirectory"])
}
