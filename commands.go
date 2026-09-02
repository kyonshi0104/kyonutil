package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "応答速度を確認します",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping": pingHandler,
}

func pingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := s.HeartbeatLatency()
	pingMs := latency.Milliseconds()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!" + fmt.Sprintf("%dms", pingMs),
		},
	})
}
