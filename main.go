package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

var webhookEndpoint string = "https://discord.com/api/webhooks/1162186214045122601/2t5gNunVxhNod_frbu8aUF05bgLWMVrGlHJ38e_qsBywE8uadla_FHxgXRQPVpQ3XDNv"

func main() {
	title := "JJK"
	episodeNumber := "690"
	imgUrl := "https://static.bunnycdn.ru/i/cache/images/c/c2/c2c8b3ae50a1b5e71d792ce9cff52431.jpg"
	watchUrl := "https://aniwave.to/watch/jujutsu-kaisen-2nd-season.ll3x3/ep-11"
	embed := makeAnimeEmbed(title, episodeNumber, imgUrl, watchUrl)
	jsonBody, _ := json.Marshal(embed)
	_, err := http.Post(webhookEndpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
}