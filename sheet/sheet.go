package sheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AntJanus/ngp-bot/config"
	"github.com/schollz/closestmatch"
)

var apiURL = "https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?key=%s"
var plusMessage = "#ERROR!"

type sheetStruct struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

type rowStruct struct {
	Name        string
	Date        string
	Platform    string
	NGP         int
	EpisodeNum  string
	EpisodeLink string
	Guest       string
	PatreonPool string
	ExactMatch  bool
	Salty       string
}

func ReadSheet(gameName string) (rowStruct, error) {
	query := fmt.Sprintf(apiURL, config.WorkBookID, config.SheetName, config.GoogleAPIKey)

	resp, err := http.Get(query)

	if err != nil {
		fmt.Println("Error is here")
		fmt.Println(err.Error())
		return rowStruct{}, err
	}

	defer resp.Body.Close()

	var sheet sheetStruct

	err = json.NewDecoder(resp.Body).Decode(&sheet)

	if err != nil {
		fmt.Println("Error is here")
		fmt.Println(err.Error())
		return rowStruct{}, err
	}

	gameTitles := []string{}
	sheetMap := make(map[string]rowStruct)

	for key, val := range sheet.Values {
		// skip first 2 rows
		if key == 0 || key == 1 {
			continue
		}

		/*
			Structure:
				Title
				Date
				Original System
				Episode #,
				NG+
		*/

		gameTitles = append(gameTitles, val[0])

		rowMap := rowStruct{
			Name:       val[0],
			Date:       checkIfExists(1, val),
			Platform:   checkIfExists(2, val),
			EpisodeNum: checkIfExists(3, val),
			NGP:        0,
			ExactMatch: false,
			Salty:      "",
		}

		if checkIfExists(4, val) != "" {
			if val[4] == "-" || val[4] == "--" {
				rowMap.NGP = -1
			} else if val[4] == "NG++" || val[4] == "NG+" {
				rowMap.NGP = 1
			}
		}

		sheetMap[strings.ToLower(val[0])] = rowMap
	}

	var gameListing rowStruct

	if val, ok := sheetMap[strings.ToLower(gameName)]; ok {
		gameListing = val
		gameListing.ExactMatch = true
	} else {
		bagSizes := []int{2}

		cm := closestmatch.New(gameTitles, bagSizes)
		gameMatch := cm.Closest(gameName)
		gameListing = sheetMap[strings.ToLower(gameMatch)]
	}

	if val, ok := config.Salty[gameListing.Name]; ok {
		gameListing.Salty = "and " + val + " is salty about it"
	}

	return gameListing, nil
}

func checkIfExists(idx int, list []string) string {
	if len(list) > idx {
		return list[idx]
	}

	return ""
}
