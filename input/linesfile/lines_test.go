package linesfile

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nuetoban/tgmailing/dto"
)

func TestGetChatsForBot(t *testing.T) {
	example := `
	0
123456
 10
 -123


	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
		return
	}

	chats, err := lf.GetChatsForBot(0)
	if err != nil {
		t.Errorf("the test met error: %v", err)
		return
	}

	expected := []dto.Chat{{ID: 0}, {ID: 123456}, {ID: 10}, {ID: -123}}
	if !cmp.Equal(expected, chats) {
		t.Errorf("GetChatsForBot returned wrong value: %v, expected: %v", chats, expected)
		return
	}
}

func TestGetChatsForBotError(t *testing.T) {
	example := `
	0
123456
 10
 -123
abc


	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
		return
	}

	_, err = lf.GetChatsForBot(0)
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
		return
	}
}

func TestGetChatsForBotErrorOnCreate(t *testing.T) {
	f, _ := os.Open("/tmp/asdlfkadsjfioasdjfoidsjf")
	_, err := New(f)
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
		return
	}
}

func TestGetBots(t *testing.T) {
	example := `
123456:foo
456789:BAR

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
		return
	}

	bots, err := lf.GetBots()
	if err != nil {
		t.Errorf("the test met error: %v", err)
		return
	}

	expected := []dto.Bot{{ID: 123456, Token: "123456:foo"}, {ID: 456789, Token: "456789:BAR"}}
	if !cmp.Equal(expected, bots) {
		t.Errorf("GetBots returned wrong value: %v, expected: %v", bots, expected)
		return
	}
}

func TestGetBotsError(t *testing.T) {
	example := `
123456:foo
456789:BAR
somestuff

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
		return
	}

	_, err = lf.GetBots()
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
		return
	}
}

func TestGetBotsErrorWrongID(t *testing.T) {
	example := `
123456:foo
456789r:BAR

	`
	lf, err := New(strings.NewReader(example))
	if err != nil {
		t.Errorf("the test met error: %v", err)
		return
	}

	_, err = lf.GetBots()
	if err == nil {
		t.Errorf("the test DID NOT meet error: %v", err)
		return
	}
}
