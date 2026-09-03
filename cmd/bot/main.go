package main

//go:generate go run ../../tools/gen_licenses.go

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	"github.com/kyonshi0104/kyonutil/internal/discord"
)

func initLogger() {
	handler := tint.NewTextHandler(os.Stderr, &tint.Options{Level: slog.LevelDebug, TimeFormat: "15:04:05"})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
func main() {
	initLogger()
	LoadEnv()
	guildID := ""

	if os.Getenv("ENVIRONMENT") == "development" {
		slog.Info("Running in development mode")
		guildID = os.Getenv("DEVELOP_GUILD_ID")
		slog.Info("Guild ID for development mode:", "guildID", guildID)
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("Error creating Discord bot:", "error", err)
		os.Exit(1)
	}

	bot.AddHandler(discord.OnReady)
	bot.AddHandler(discord.OnMessageCreate)
	bot.AddHandler(discord.OnInteractionCreate)
	if err := bot.Open(); err != nil {
		slog.Error("Error opening Discord bot:", "error", err)
		os.Exit(1)
	}

	defer bot.Close()
	slog.Info("Bot is now running. Press Ctrl+C to exit.")

	go startPingStatusUpdater(bot)

	_, err = bot.ApplicationCommandBulkOverwrite(bot.State.User.ID, guildID, discord.Commands)

	if err != nil {
		slog.Error("Error overwriting application commands:", "error", err)
		os.Exit(1)
	} else {
		slog.Info("Application commands registered successfully.")
	}

	// Wait for a signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}

func startPingStatusUpdater(s *discordgo.Session) {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		<-ticker.C
		ping := s.HeartbeatLatency().Milliseconds()
		statusText := fmt.Sprintf("Ping: %dms", ping)

		err := s.UpdateCustomStatus(statusText)
		if err != nil {
			slog.Warn("Failed to update status", "error", err)
		}
	}
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Warning: .env file not found. Relying on system environment variables")
	}
}
