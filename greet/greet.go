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
 *   greet help
 *   greet --lang spanish --age 19 xxxx
 *   greet template help
 */
func main() {
	var language string //go变量申明，同种类型可以放在一行，比如var a, b string，不同类型需要单独行
	var age int

	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "fight the loneliness!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang, l", //逗号分隔，后面的表示短标签
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
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`", //反引号圈引表示这是一个占位符，结果为：--config FILE, -c FILE   Load configuration from FILE
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:     "template",
			Aliases:  []string{"t"},
			Usage:    "options for task templates",
			Category: "Template actions",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
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
