package linesfile

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/nuetoban/tgmailing/dto"
)

type Src struct {
	lines []string
}

func New(file io.Reader) (*Src, error) {
	var src Src

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return &src, err
	}

	src.lines = strings.Fields(string(b))

	return &src, nil
}

func (s *Src) GetChatsForBot(_ int) ([]dto.Chat, error) {
	var chats []dto.Chat

	for _, c := range s.lines {
		chatInt, err := strconv.Atoi(c)
		if err != nil {
			return chats, err
		}
		chats = append(chats, dto.Chat{ID: int64(chatInt)})
	}

	return chats, nil
}

func (s *Src) GetBots() ([]dto.Bot, error) {
	var bots []dto.Bot

	for _, v := range s.lines {
		botIDStr := strings.Split(v, ":")
		if len(botIDStr) < 2 {
			return bots, fmt.Errorf("all the lines should contain ':'")
		}

		botID, err := strconv.Atoi(botIDStr[0])
		if err != nil {
			return bots, err
		}

		bots = append(bots, dto.Bot{ID: botID, Token: v})
	}

	return bots, nil
}
