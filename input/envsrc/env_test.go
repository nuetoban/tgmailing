package envsrc

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nuetoban/tgmailing/dto"
)

func TestGetBots(t *testing.T) {
	os.Setenv("TEST_ENV_NAME", "123456:aaaaaa")

	e, err := New("TEST_ENV_NAME")
	if err != nil {
		t.Errorf("test returned the error: %v", err)
		return
	}

	bots, err := e.GetBots()
	if err != nil {
		t.Errorf("GetBots() returned the error: %v", err)
		return
	}

	expected := []dto.Bot{{ID: 123456, Token: "123456:aaaaaa"}}
	if !cmp.Equal(bots, expected) {
		t.Errorf("GetBots() wrong value: %v, expected: %v", bots, expected)
		return
	}
}

func TestGetBotsErrorInvalidFormat(t *testing.T) {
	// No ":"
	os.Setenv("TEST_ENV_NAME", "123456aaaaaa")

	_, err := New("TEST_ENV_NAME")
	if err == nil {
		t.Errorf("test DID NOT return the error")
		return
	}
}

func TestGetBotsErrorInvalidID(t *testing.T) {
	// No ":"
	os.Setenv("TEST_ENV_NAME", "123456b:aaaaaa")

	_, err := New("TEST_ENV_NAME")
	if err == nil {
		t.Errorf("test DID NOT return the error")
		return
	}
}
