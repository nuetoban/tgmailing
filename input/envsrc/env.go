package envsrc

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nuetoban/tgmailing/dto"
)

type Src struct {
	b []dto.Bot
}

func New(envName string) (*Src, error) {
	var src Src

	b := os.Getenv(envName)
	if b == "" {
		return &src, fmt.Errorf("provided env %s is empty", envName)
	}

	botIDStr := strings.Split(b, ":")[0]
	botID, err := strconv.Atoi(botIDStr)
	if err != nil {
		return &src, err
	}

	src.b = append(src.b, dto.Bot{ID: botID, Token: b})

	return &src, nil
}

func (s *Src) GetBots() ([]dto.Bot, error) {
	return s.b, nil
}
