package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/urfave/cli"
)

/*
 * 在greet目录中执行go install，会在$GOPATH/bin目录中安装greet命令
 * 使用方法：
 *   greet --lang spanish --age 19 xxxx
 */
func main() {
	var language string //go变量申明，同种类型可以放在一行，比如var a, b string，不同类型需要单独行
	var age int

	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "fight the loneliness!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang",
			Value:       "english", //设置默认值的参数
			Usage:       "language for the greeting",
			Destination: &language,
		},
		cli.IntFlag{
			Name:        "age",
			Value:       18,
			Usage:       "guy's age",
			Destination: &age,
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println(reflect.TypeOf(c.Args())) //返回的Args类型，基础类型是字符串切片 type Args []string

		name := "someone"
		if c.NArg() > 0 { //有几个参数
			name = c.Args()[0]
		}
		if language == "spanish" {
			fmt.Println("Hola", name)
		} else {
			fmt.Println("Hello", name)
		}
		if age == 19 {
			fmt.Println("Hello", age)
		}
		return nil
	}
	app.Run(os.Args)
}
