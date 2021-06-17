package jsonsrc

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nuetoban/tgmailing/dto"
)

func TestGetChatsForBot(t *testing.T) {
	example := `
[
	{"id": 0},
	{"id": 123456},

	{"id": 10},
	{"id": -123}
]

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	chats, err := lf.GetChatsForBot(0)
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	expected := []dto.Chat{{ID: 0}, {ID: 123456}, {ID: 10}, {ID: -123}}
	if !cmp.Equal(expected, chats) {
		t.Errorf("GetChatsForBot returned wrong value: %v, expected: %v", chats, expected)
	}
}

func TestGetChatsForBotErrorWrongChatID(t *testing.T) {
	example := `
[
	{"id": 0},
	{"id": 123456},

	{"id": "abc"},
	{"id": -123}
]

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	_, err = lf.GetChatsForBot(0)
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
	}
}

func TestGetChatsForBotErrorInvalidJSON(t *testing.T) {
	example := `
[
	{"id": 0},
	{"id": 123456},

	{"id": "abc"},
	{"id": -123}
]
// Make it invalid by this command

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	_, err = lf.GetChatsForBot(0)
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
	}
}

func TestGetChatsForBotErrorOnCreate(t *testing.T) {
	f, _ := os.Open("/tmp/asdlfkadsjfioasdjfoidsjf")
	_, err := New(f)
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
	}
}

func TestGetBots(t *testing.T) {
	example := `
[
	{"id": 123456, "token": "123456:foo"},
	{"id": 456789, "token": "456789:BAR"}
]

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	bots, err := lf.GetBots()
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	expected := []dto.Bot{{ID: 123456, Token: "123456:foo"}, {ID: 456789, Token: "456789:BAR"}}
	if !cmp.Equal(expected, bots) {
		t.Errorf("GetBots returned wrong value: %v, expected: %v", bots, expected)
	}
}

func TestGetBotsErrorInvalidJSON(t *testing.T) {
	example := `
[
	{"id": 123456, "token": "123456:foo"},
	{"id": 456789, "token": "456789:BAR"},
]

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	_, err = lf.GetBots()
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
	}
}

func TestGetBotsErrorWrongID(t *testing.T) {
	example := `
[
	{"id": 123456, "token": "123456:foo"},
	{"id": "456789a", "token": "456789:BAR"},
]

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
		return
	}

	_, err = lf.GetBots()
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
	}
}

func TestGetScheduledAd(t *testing.T) {
	example := `
	{
		"message": {
			"message_type": "text",
			"interval": 0.1,
			"test": false,
			"text": "poebat"
		}
	}

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	sad, err := lf.GetScheduledAd(0, 0, 0)
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	expected := dto.ScheduledAd{Message: dto.Ad{MessageType: "text", Interval: 0.1, Test: false, Text: "poebat"}}
	if !cmp.Equal(expected, sad) {
		t.Errorf("GetBots returned wrong value: %v, expected: %v", sad, expected)
	}
}

func TestGetScheduledAdErrorInvalidJSON(t *testing.T) {
	example := `
	{
		"message": {
			"message_type": "text",
			"interval": 0.1,
			"test": false,
			"text": "poebat"
		},,,
	}

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
	}

	_, err = lf.GetScheduledAd(0, 0, 0)
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
		return
	}
}
