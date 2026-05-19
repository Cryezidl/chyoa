package cyoa

import (
	"encoding/json"
	"log/slog"
	"os"
)

type Chapter struct {
	Title      string    `json:"title"`
	Paragraphs []string  `json:"story"`
	Options    []Options `json:"options"`
}

type Options struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Story map[string]Chapter

func LoadStory(fileName string, logger *slog.Logger) (Story, error) {
	//Открыть файл
	file, err := os.Open(fileName)
	if err != nil {

		return nil, err
	}
	defer file.Close()
	//Создать Story
	var s Story
	//Декодировать json в Story
	if err := json.NewDecoder(file).Decode(&s); err != nil {

		return nil, err
	}
	//Вернуть Story
	return s, nil
}

func GetChapter(s Story, storyName string) (Chapter, bool) {
	val, ok := s[storyName]
	return val, ok
}
