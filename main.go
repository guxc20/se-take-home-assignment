package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var logChannel chan string
func main() {
	logChannel = make(chan string, 100)

	// 启动日志打印 goroutine
	go func() {
		for msg := range logChannel {
			fmt.Println(msg)
			fmt.Print("> ")
		}
	}()

	controller := NewController()
	reader := bufio.NewReader(os.Stdin)

	logChannel <- fmt.Sprintf("=== 麦当劳烹饪机器人系统 ===")
	logChannel <- fmt.Sprintf("指令: new normal (nn) | new vip (nv) | +bot (+b) | -bot (-b) | status (st) | exit (ex)")

	aliasMap := map[string]string{
		"nv": "new vip",
		"nn": "new normal",
		"+b": "+bot",
		"-b": "-bot",
		"st":  "status",
		"ex":  "exit",
	}

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// 替换别名
		if fullCmd, ok := aliasMap[input]; ok {
			input = fullCmd
		}

		switch input {
		case "new normal":
			controller.CreateOrder(false)
		case "new vip":
			controller.CreateOrder(true)
		case "+bot":
			controller.AddRobot()
		case "-bot":
			controller.RemoveRobot()
		case "status":
			controller.PrintStatus()
		case "exit":
			logChannel <- fmt.Sprintf("退出程序")
			return
		default:
			logChannel <- fmt.Sprintf("未知命令\n")
		}
	}
}
