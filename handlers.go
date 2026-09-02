package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "テストなう")
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
}

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name

	if handler, ok := commandHandlers[commandName]; ok {
		handler(s, i)
	} else {
		fmt.Printf("Unhandled command: %s\n", commandName)
	}
}
