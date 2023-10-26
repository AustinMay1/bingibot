package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Tile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	client          = &http.Client{}
	intOptionMinVal = 1.0

	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "setup",
			Description: "setup bingo team channels.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "num-of-teams",
					Description: "number of team channels to create",
					MinValue:    &intOptionMinVal,
					MaxValue:    99,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "tiles-list",
					Description: "upload the JSON file containing the bingo tiles",
					Required:    true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"setup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			start := time.Now()
			// create a deferred response while bot works
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})

			var tiles []*Tile
			var numOfTeams int
			cmdData := i.ApplicationCommandData().Options
			configFile := i.ApplicationCommandData().Resolved.Attachments
			options := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(cmdData))

			for _, opt := range cmdData {
				options[opt.Name] = opt
			}

			if option, ok := options["num-of-teams"]; ok {
				numOfTeams = int(option.IntValue())
			} // extract value from cmdData. Seems unecessarily complex.

			for _, file := range configFile {
				req, err := http.NewRequest("GET", file.URL, nil) // after uploading file to cdn.discord, fetch it
				req.Header.Add("Content-Type", "application/json")

				if err != nil {
					log.Printf("Error: %v", err)
				}

				res, err := client.Do(req)

				if err != nil {
					log.Printf("Error: %v", err)
				}

				body, err := io.ReadAll(res.Body)

				if err != nil {
					log.Printf("Error: %v", err)
				}

				err = json.Unmarshal(body, &tiles)

				if err != nil {
					log.Printf("Error: %v", err)
				}

				for _, tile := range tiles {
					fmt.Printf("\n%v\n", tile)
				}
			} // are all these err checks really necessary?

			for j := 1; j <= numOfTeams; j++ {
				go func(name int) {
					chName := fmt.Sprintf("team-%v", name)
					ch, err := s.GuildChannelCreate(i.GuildID, chName, discordgo.ChannelTypeGuildForum)

					if err != nil {
						log.Fatalf("Failed to create channel: %v", err)
					}

					log.Printf("%v\n", ch.Name)

					for _, tile := range tiles {
						thread, err := s.ForumThreadStart(ch.ID, tile.Name, 0, tile.Description)

						if err != nil {
							log.Fatalf("Failed to create thread: %v", thread.ID)
						}
					}
				}(j)
			}
			// create interaction response and update deferred response
			deferredMsg := fmt.Sprintf("teams created: %v", numOfTeams)

			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &deferredMsg,
			})
			elapsed := time.Since(start).Seconds()
			fmt.Printf("\n%.2fs\n", elapsed)
		},
	}
)
