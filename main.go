package main

import (
	"fmt"
	"github.com/AntJanus/ngp-bot/bot"
	"github.com/AntJanus/ngp-bot/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
