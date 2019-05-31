package main

import (
	"fmt"

	"git.llsapp.com/common/dva-go/dva"
)

func main() {
	conf := &dva.Config{App: "telisgo", Env: "production"}
	cli, err := dva.NewClient(conf, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	config, err := cli.Get("test_key")
	fmt.Println(config, err)
	secret, err := cli.GetSecret("db_data_source")
	fmt.Println(secret, err)
}
