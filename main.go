package main

import (
	"muxCli/cmd"
	"muxCli/config"
)

func main() {
	// 读配置文件
	config.ReadConfig()

	err := cmd.RootCmd.Execute()
	if err != nil {
		return
	}
}
