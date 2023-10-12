package commands

import "github.com/bwmarrin/discordgo"

var (
	Commands = []*discordgo.ApplicationCommand {
		{
			Name: "create",
			Description: "create a channel",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		"create": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: "testing the /create command!",
				},
			})
		},
	}
)