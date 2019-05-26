package sheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AntJanus/budget-arcade-bot/config"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/schollz/closestmatch"
)

var apiURL = "https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?key=%s"

type sheetStruct struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

/*
RowStruct describes structure of spreadsheet data
*/
type RowStruct struct {
	Name        string
	Date        string
	Platform    string
	Approval    int
	EpisodeNum  string
	EpisodeLink string
	Guest       string
	PatreonPool string
	ExactMatch  bool
	Salty       string
}

/*
ReadSheet fetches spreadsheet, reads it, and matches a game to input gameName
*/
func ReadSheet(gameName string) (RowStruct, error) {
	query := fmt.Sprintf(apiURL, config.WorkBookID, config.SheetName, config.GoogleAPIKey)

	resp, err := http.Get(query)

	if err != nil {
		fmt.Println("Error is here")
		fmt.Println(err.Error())
		return RowStruct{}, err
	}

	defer resp.Body.Close()

	var sheet sheetStruct

	err = json.NewDecoder(resp.Body).Decode(&sheet)

	if err != nil {
		fmt.Println("Error is here")
		fmt.Println(err.Error())
		return RowStruct{}, err
	}

	gameTitles := []string{}
	sheetMap := make(map[string]RowStruct)

	for key, val := range sheet.Values {
		// skip first 2 rows
		if key == 0 {
			continue
		}

		/*
			Structure:
				Name
				Platform
				Episode
				Seal of Approval
		*/

		gameTitles = append(gameTitles, val[0])

		rowMap := RowStruct{
			Name:       val[0],
			Platform:   checkIfExists(1, val),
			EpisodeNum: checkIfExists(2, val),
			Approval:   0,
			ExactMatch: false,
			Salty:      "",
		}

		if checkIfExists(3, val) != "" {
			if val[3] == "No" || val[3] == "no" {
				rowMap.Approval = -1
			} else if val[3] == "Yes" || val[3] == "yes" {
				rowMap.Approval = 1
			}
		}

		sheetMap[strings.ToLower(val[0])] = rowMap
	}

	var gameListing RowStruct

	if val, ok := sheetMap[strings.ToLower(gameName)]; ok {
		gameListing = val
		gameListing.ExactMatch = true
	} else {
		bagSizes := []int{2, 3, 4}

		cm := closestmatch.New(gameTitles, bagSizes)
		gameMatch := cm.Closest(gameName)
		distance := fuzzy.LevenshteinDistance(gameName, gameMatch)

		if distance < 12 {
			gameListing = sheetMap[strings.ToLower(gameMatch)]
		}
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
