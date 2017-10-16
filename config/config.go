package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Token        string
	BotPrefix    string
	MashapeKey   string
	GoogleAPIKey string
	WorkBookID   string
	SheetName    string
	Salty        map[string]string

	config *configStruct
)

type configStruct struct {
	Token        string            `json:"Token"`
	BotPrefix    string            `json:"BotPrefix"`
	MashapeKey   string            `json:"MashapeKey"`
	GoogleAPIKey string            `json:"GoogleAPIKey"`
	WorkBookID   string            `json:"WorkBookID"`
	SheetName    string            `json:"SheetName"`
	Salty        map[string]string `json:"Salty"`
}

func ReadConfig() error {
	fmt.Println("Reading from config file...")

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())

		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())

		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix
	MashapeKey = config.MashapeKey
	GoogleAPIKey = config.GoogleAPIKey
	WorkBookID = config.WorkBookID
	SheetName = config.SheetName
	Salty = config.Salty

	return nil
}
