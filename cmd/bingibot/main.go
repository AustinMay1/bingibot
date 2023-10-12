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
        log.Fatalf("Error reading config file.")
    }

    token, ok := viper.Get("DISCORD_TOKEN").(string)

    if !ok {
        log.Fatalf("Invalid token.")
    }

    discord, err := discordgo.New("Bot " + token)

    if err != nil {
        log.Fatal("Invalid authentication.")
    }

    discord.AddHandler(pingpong)
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
            log.Panicf("Cannot create '%v'", v.Name)
        }
        registeredCommands[i] = cmd
    }

    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <- sc

    discord.Close()
}

func pingpong(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    if m.Content == "ping" {
        s.ChannelMessageSend(m.ChannelID, "pong!")
    }

    if m.Content == "pong" {
        s.ChannelMessageSend(m.ChannelID, "ping!")
    }
}
