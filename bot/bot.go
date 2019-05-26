package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/AntJanus/budget-arcade-bot/config"
	"github.com/AntJanus/budget-arcade-bot/igdb"
	"github.com/AntJanus/budget-arcade-bot/sheet"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *discordgo.Session

var (
	yesEmoji   = "white_check_mark"
	noEmoji    = "no_entry"
	plusEmoji  = "white_check_mark"
	minusEmoji = "no_entry"
)

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
			query := strings.TrimPrefix(command, "check ")
			game, err := igdb.Search(query)

			if err != nil {
				fmt.Println("Error is here")
				fmt.Println(err.Error())
				return
			}

			if game.Name == "" {
				message := fmt.Sprintf("Cannot find Game: %s", query)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)

				return
			}

			releaseDate := game.FirstReleaseDate
			unixDate := time.Unix(releaseDate, 0)
			humanDate := unixDate.Format("01/02/2006")

			// cover
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", game.URL))

			// game information
			message := ""

			if releaseDate == 0 {
				message = fmt.Sprintf("Game: %s \nDate: [unknown] \n", game.Name)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)
			} else {
				message = fmt.Sprintf("Game: %s \nDate: %s", game.Name, humanDate)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)
			}

			query = game.Name
			gameMatch, _ := sheet.ReadSheet(query)
			approvalStatus := ""
			message = ""

			if len(gameMatch.Name) == 0 {
				message += "Could not find game in master list"
			} else {
				if gameMatch.ExactMatch == false {
					message += "The closest match I could find: \n"
					message += fmt.Sprintf("Game: %s\n", gameMatch.Name)
				} else {
					message += "Game is in the master list\n"
				}
			}

			if gameMatch.Approval == 1 {
				approvalStatus = plusEmoji
				message += fmt.Sprintf("Approval: :%s: %s\n", approvalStatus, gameMatch.Salty)
				message += fmt.Sprintf("Ep#: %s", gameMatch.EpisodeNum)
			} else if gameMatch.Approval == -1 {
				approvalStatus = minusEmoji
				message += fmt.Sprintf("Approval: :%s: %s", approvalStatus, gameMatch.Salty)
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}
	}

}
