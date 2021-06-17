package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"database/sql/driver"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	ID    int    `json:"id" db:"id"`
	Token string `json:"token" db:"token"`
}

type Chat struct {
	ID int64 `json:"id" db:"id"`
}

type Ad struct {
	MessageType string  `json:"message_type"`
	Interval    float64 `json:"interval"` // Default: 0.1
	Test        bool    `json:"test"`

	Text        string                   `json:"text"`                // Will be translated to "caption" if the type is not TEXT
	FileBlob    *string                  `json:"file_blob,omitempty"` // used for PHOTO, ANIMATION and VIDEO
	ReplyMarkup *tb.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func NewAd() Ad {
	a := Ad{}
	a.Interval = 0.1
	return a
}

func (a Ad) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Ad) Scan(value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("wrong type")
	}
	return json.Unmarshal(v, a)
}

type ScheduledAd struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`

	Message Ad `json:"message"`
}

type Statistics struct {
	BotID                  int       `json:"bot_id"`
	BotName                string    `json:"bot_name"`
	SuccessfulSendAttempts int       `json:"successful_send_attempts"`
	FailedSendAttempts     int       `json:"failed_send_attemplts"`
	StartTime              time.Time `json:"start_time"`
}

type MultipleBotsStatistics struct {
	ID         int64
	Statistics []Statistics
	StartTime  time.Time
}
