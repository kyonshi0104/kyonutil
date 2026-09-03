package discord

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/kyonshi0104/kyonutil/internal/tld"
)

//go:embed licenses.csv
var licenseCSV string

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "応答速度を確認します",
	},
	{
		Name:        "charcount",
		Description: "文字数をカウントします",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "カウントする文字列を入力してください",
				Required:    true,
			},
		},
	},
	{
		Name:        "tldcheck",
		Description: "TLDの登録状況を確認します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "domain",
				Description: "確認するドメイン名を入力してください",
				Required:    true,
			},
		},
	},
	{
		Name:        "license",
		Description: "使用しているライブラリのライセンス情報を表示します",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":      pingHandler,
	"charcount": charCountHandler,
	"tldcheck":  tldCheckHandler,
	"license":   licenseHandler,
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

func charCountHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	start := time.Now()
	text := i.ApplicationCommandData().Options[0].StringValue()
	count := len([]rune(text))
	textWithoutSpaces := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1 // -1を返すとその文字は削除される
		}
		return r
	}, text)

	countWithoutSpaces := len([]rune(textWithoutSpaces))

	elapsed := time.Since(start).String()

	embed := &discordgo.MessageEmbed{
		Title: "文字数",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "総文字数",
				Value:  fmt.Sprintf("%d", count),
				Inline: true,
			},
			{
				Name:   "文字数（スペース除く）",
				Value:  fmt.Sprintf("%d", countWithoutSpaces),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Processing time: %s", elapsed),
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func tldCheckHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	start := time.Now()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	domain := i.ApplicationCommandData().Options[0].StringValue()
	results := tld.CheckCommonTLDs(domain)

	var tlds []string
	for k := range results {
		tlds = append(tlds, k)
	}
	sort.Strings(tlds)

	description := ""
	for tld, available := range results {
		status := "✅ Available"
		if !available {
			status = "❌ Unavailable"
		}
		description += fmt.Sprintf("**%s**: %s\n", tld, status)
	}

	elapsed := time.Since(start).String()

	embed := &discordgo.MessageEmbed{
		Title:       domain,
		Color:       0x00ff00,
		Description: description,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Processing time: %s", elapsed),
		},
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}

func licenseHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lines := strings.Split(strings.TrimSpace(licenseCSV), "\n")

	var descBuilder strings.Builder
	descBuilder.WriteString("KyonUtil is built with the following open source software:\n\n")

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) >= 3 {
			module := parts[0]
			url := parts[1]
			licenseName := parts[2]

			descBuilder.WriteString(fmt.Sprintf("・[%s](%s) - **%s**\n", module, url, licenseName))
		}
	}

	finalDesc := descBuilder.String()
	if len(finalDesc) > 4000 {
		finalDesc = finalDesc[:4000] + "...\n(truncated)"
	}

	// Embedを生成して返信
	embed := &discordgo.MessageEmbed{
		Title:       "Licenses of Libraries Used in KyonUtil",
		Description: finalDesc,
		Color:       0x0099ff,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
