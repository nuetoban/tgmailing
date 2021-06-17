package jsonsrc

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/nuetoban/tgmailing/dto"
)

type Src struct {
	content []byte
}

func New(r io.Reader) (*Src, error) {
	var (
		src Src
		err error
	)

	src.content, err = ioutil.ReadAll(r)
	if err != nil {
		return &src, err
	}

	return &src, nil
}

func (s *Src) GetScheduledAd(_ int, _ int, _ int) (dto.ScheduledAd, error) {
	var sad dto.ScheduledAd

	err := json.Unmarshal(s.content, &sad)
	if err != nil {
		return sad, err
	}
	return sad, nil
}

func (s *Src) GetChatsForBot(_ int) ([]dto.Chat, error) {
	var chats []dto.Chat

	err := json.Unmarshal(s.content, &chats)
	if err != nil {
		return chats, err
	}

	return chats, nil
}

func (s *Src) GetBots() ([]dto.Bot, error) {
	var bots []dto.Bot

	err := json.Unmarshal(s.content, &bots)
	if err != nil {
		return bots, err
	}

	return bots, nil
}
