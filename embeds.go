package main

import "fmt"

type Thumbnail struct {
	Url string `json:"url"`
}

type Embed struct {
	Type string `json:"type"`
	Title string `json:"title"`
	Description string `json:"description"`
	Color int `json:"color"`
	Thumbnail Thumbnail `json:"thumbnail"`
	Url string `json:"url"`
}

type AnimeEmbed struct {
	Content string `json:"content"`
	Tts bool `json:"tts"`
	Embeds []Embed `json:"embeds"`
}

func makeAnimeEmbed(title, episodeNumber, imgUrl, watchUrl string) AnimeEmbed {
	embed := Embed{"rich", "New Episode Detected", fmt.Sprintf("%s Episode %s has just released!", title, episodeNumber), 0x00FFFF, Thumbnail{imgUrl}, watchUrl}
	return AnimeEmbed{"", false, []Embed{embed}}
}