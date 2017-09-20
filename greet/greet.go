package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

/*
 * 在greet目录中执行go install，会在$GOPATH/bin目录中安装greet命令
 */
func main() {
	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "fight the loneliness!"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello friend!")
		return nil
	}
	app.Run(os.Args)
}
