package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AustinMay1/bingibot/pkg/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Error reading config file.")
	}

	token, ok := viper.Get("DISCORD_TOKEN").(string)

	if !ok {
		log.Fatal("Invalid key.")
	}

	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatal("Invalid token.")
	}
	
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()

	if err != nil {
		log.Fatal("Error opening connection")
	}

	fmt.Printf("Logged in as: %v#%v", discord.State.User.Username, discord.State.User.Discriminator)

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.Commands))
	for i, v := range commands.Commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v': %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-kill

	discord.Close()
}
