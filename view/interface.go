package view

import (
	"encoding/json"

	"github.com/nuetoban/tgmailing/dto"
)

type StatisticsProvider interface {
	Statistics() dto.MultipleBotsStatistics
}

type Error struct {
	Error error
}

func (e Error) JSON() []byte {
	content, _ := json.Marshal(e)
	return content
}
