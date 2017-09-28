package bot

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
  "net/url"
	"ngp-bot/config"
	"strings"
)

var BotID string
var goBot *discordgo.Session

var apiURL = "https://igdbcom-internet-game-database-v1.p.mashape.com/games/?fields=name%2Crelease_dates&limit=1&offset=0&search="
var client = &http.Client{}

type GameStruct []struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ReleaseDates []struct {
		Category int    `json:"category"`
		Platform int    `json:"platform"`
		Date     int64  `json:"date"`
		Human    string `json:"human"`
		Y        int    `json:"y"`
		M        int    `json:"m"`
	} `json:"release_dates"`
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	return
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		command := strings.TrimPrefix(m.Content, config.BotPrefix)

		fmt.Println("Runs")
		fmt.Println(command)

		if m.Author.ID == BotID {
			return
		}

		if strings.HasPrefix(command, "check") {
			fmt.Println("Checking")
			query := strings.TrimPrefix(command, "check ")
      urlQuery := url.QueryEscape(query)
			req, err := http.NewRequest("GET", apiURL+urlQuery, nil)

			req.Header.Add("X-Mashape-Key", config.MashapeKey)
			req.Header.Add("Accepts", "application/json")

			resp, err := client.Do(req)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			defer resp.Body.Close()

			var games GameStruct

			err = json.NewDecoder(resp.Body).Decode(&games)

			if err != nil {
				fmt.Println("Error is here")
				fmt.Println(err.Error())
				return
			}

			fmt.Println("Done")

			message := fmt.Sprintf("Game: %s \nDate: %s", games[0].Name, games[0].ReleaseDates[0].Human)

			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}
	}

}
