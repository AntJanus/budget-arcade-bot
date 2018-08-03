package igdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/AntJanus/ngp-bot/config"
)

var (
	apiURL       = "https://api-2445582011268.apicast.io"
	searchURL    = apiURL + "/games/?search="
	searchFields = "id,name,first_release_date,url,summary,rating,time_to_beat,cover"
)

var client = &http.Client{}

type GameStruct struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	FirstReleaseDate int64  `json:"first_release_date,omitempty"`
	Cover            struct {
		URL          string `json:"url"`
		CloudinaryID string `json:"cloudinary_id"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
	} `json:"cover,omitempty"`
}
type GameList []GameStruct

func Search(gameName string) (GameStruct, error) {
	urlQuery := url.QueryEscape(gameName)
	req, err := http.NewRequest("GET", searchURL+urlQuery+"&fields="+searchFields, nil)

	req.Header.Add("user-key", config.IGDBKey)
	req.Header.Add("Accepts", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		var empty GameStruct

		return empty, err
	}

	defer resp.Body.Close()

	var games GameList

	err = json.NewDecoder(resp.Body).Decode(&games)

	if err != nil {
		fmt.Println(err.Error())
		var empty GameStruct

		return empty, err
	}

	game := games[0]

	if strings.HasPrefix(game.Cover.URL, "//") {
		game.Cover.URL = "https:" + game.Cover.URL
	}

	return game, err
}
