package handlers

import (
	"strconv"

	"github.com/darenliang/jikan-go"
)

type AnimeHandler struct {
}

func MakeAnimeHandler() *AnimeHandler {
	return &AnimeHandler{}
}

func (a AnimeHandler) QueryAnimeData (malId string) *jikan.AnimeById {
	id, _ := strconv.Atoi(malId)
	anime, err := jikan.GetAnimeById(id)
	if err != nil {
		return nil
	}
	return anime
}
