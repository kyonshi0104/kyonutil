package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

func initLogger() {
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05",
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
func main() {
	initLogger()
	loadEnv()
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

	bot.AddHandler(onReady)
	bot.AddHandler(onMessageCreate)
	bot.AddHandler(onInteractionCreate)
	if err := bot.Open(); err != nil {
		slog.Error("Error opening Discord bot:", "error", err)
		os.Exit(1)
	}

	defer bot.Close()
	slog.Info("Bot is now running. Press Ctrl+C to exit.")

	_, err = bot.ApplicationCommandBulkOverwrite(bot.State.User.ID, guildID, commands)

	if err != nil {
		slog.Error("Error overwriting application commands:", "error", err)
		os.Exit(1)
	}

	// Wait for a signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Warning: .env file not found. Relying on system environment variables")
	}
}
