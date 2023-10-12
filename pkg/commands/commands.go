package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type tile struct {
	name string
	archiveDuration int
	content string
}

var (
	Tiles = []*tile{
		{
			name: "500 Herbi KC",
			archiveDuration: 0,
			content: "Get 500 Herbibore KC.",
		},
		{
			name: "1k Tithe Farm Points",
			archiveDuration: 0,
			content: "Get 1000 points in the Tithe Farm minigame.",
		},
	}

	Commands = []*discordgo.ApplicationCommand {
		{
			Name: "create",
			Description: "create a channel",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		"create": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmt.Printf("\n/create [Guild: %v]", i.GuildID)
			ch, err := s.GuildChannelCreate(i.GuildID, "test", discordgo.ChannelTypeGuildForum)

			if err != nil {
				log.Fatalf("Failed to create channel: %v", err.Error())
			}

			for _, v := range Tiles {
				_, err := s.ForumThreadStart(ch.ID, v.name, v.archiveDuration, v.content)

				if err != nil {
					log.Fatalf("Failed to create post %v in <#%v>", v.name, ch.ID)
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: fmt.Sprintf("Channel <#%v> created.", ch.ID),
				},
			})
		},
	}
)