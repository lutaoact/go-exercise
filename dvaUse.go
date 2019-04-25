package main

import (
	"fmt"

	"git.llsapp.com/common/dva-go/dva"
)

func main() {
	conf := &dva.Config{App: "telisgo"}
	cli, err := dva.NewClient(conf, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	config, err := cli.Get("test_config")
	fmt.Println(config, err)
	secret, err := cli.GetSecret("redis_password")
	fmt.Println(secret, err)
}
