package discord

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func OnReady(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateCustomStatus("Loading...")
	slog.Info("Bot is ready.", "username", event.User.Username, "discriminator", event.User.Discriminator)
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
}

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
