package util

import (
	"errors"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SuccessEmbed(message string) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Description: message,
		Color:       ColorGreen,
	}
	return &embed
}

func ErrorEmbed(message string) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Description: message,
		Color:       ColorRed,
	}
	return &embed
}

func ParseMalIdFromUrl(malUrl string) (string, error) {
	u, err := url.Parse(malUrl)
	if err != nil {
		return "", err
	}
	escaped := u.EscapedPath()
	base := "/anime/"
	if !strings.Contains(escaped, base) {
		return "", errors.New("could not parse MyAnimeList ID from URL")
	}
	base_cut := escaped[len(base):]
	end_id := strings.Index(base_cut, "/")
	return base_cut[:end_id], nil
}
