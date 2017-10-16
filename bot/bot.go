package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/AntJanus/ngp-bot/config"
	"github.com/AntJanus/ngp-bot/sheet"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *discordgo.Session

var apiURL = "https://igdbcom-internet-game-database-v1.p.mashape.com/games/?fields=name%2Crelease_dates&limit=1&offset=0&search="
var client = &http.Client{}

var yesEmoji = "white_check_mark"
var noEmoji = "no_entry"
var plusEmoji = "white_check_mark"
var minusEmoji = "no_entry"

type ReleaseDate struct {
	Category int    `json:"category"`
	Platform int    `json:"platform"`
	Date     int64  `json:"date"`
	Human    string `json:"human"`
	Y        int    `json:"y"`
	M        int    `json:"m"`
}

type ReleaseDates []ReleaseDate

type GameStruct []struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	ReleaseDates ReleaseDates `json:"release_dates"`
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

func (r ReleaseDates) Less(i, j int) bool { return r[i].Date < r[j].Date }
func (r ReleaseDates) Len() int           { return len(r) }
func (r ReleaseDates) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		command := strings.TrimPrefix(m.Content, config.BotPrefix)

		fmt.Println("Runs")
		fmt.Println(command)

		if m.Author.ID == BotID {
			return
		}

		if strings.HasPrefix(command, "status") {
			fmt.Println("Checking...")
			query := strings.TrimPrefix(command, "status ")
			gameMatch, _ := sheet.ReadSheet(query)
			ngpStatus := ""
			message := ""

			if gameMatch.NGP == 1 {
				ngpStatus = plusEmoji
			} else if gameMatch.NGP == -1 {
				ngpStatus = minusEmoji
			}

			if gameMatch.ExactMatch == false {
				message += "The closest match I could find: \n"
			}

			message += fmt.Sprintf("Game: %s \nDate: %s \nNGP: :%s: %s", gameMatch.Name, gameMatch.Date, ngpStatus, gameMatch.Salty)

			if gameMatch.NGP != 0 {
				message += fmt.Sprintf("\nEpisode Number: %s\nEpisode Link: %s", gameMatch.EpisodeNum, gameMatch.EpisodeLink)
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, message)

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

			if len(games) == 0 {
				message := fmt.Sprintf("Cannot find Game: %s", query)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)

			} else {

				fmt.Println("Done")

				game := games[0]
				releaseDates := game.ReleaseDates
				sort.Sort(ReleaseDates(releaseDates))

				humanDate := ""

				if len(releaseDates) > 0 {
					humanDate = releaseDates[0].Human
				} else {
					humanDate = "No available date"
				}

				nowDate := time.Now()
				fmt.Println(releaseDates[0].Date)
				unixDate := time.Unix(releaseDates[0].Date/1000, 0)
				dateDifference := nowDate.Sub(unixDate)
				yearDiff := dateDifference.Hours() / 24 / 365

				statusEmoji := ""

				if yearDiff > 15 {
					statusEmoji = yesEmoji
				} else {
					statusEmoji = noEmoji
				}

				message := fmt.Sprintf("Game: %s \nDate: %s \nEligible: :%s:", game.Name, humanDate, statusEmoji)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)
			}
		}
	}

}
