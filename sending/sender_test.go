package sending

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nuetoban/tgmailing/dto"

	tb "gopkg.in/tucnak/telebot.v2"
)

type dataProvider struct {
	ad dto.ScheduledAd
}

func (d *dataProvider) GetBots() ([]dto.Bot, error) {
	var m []dto.Bot

	botToken := os.Getenv("SENDER_TOKEN")

	tokenSplit := strings.Split(botToken, ":")
	if len(tokenSplit) != 2 {
		return m, fmt.Errorf("invalid token")
	}

	botID, err := strconv.Atoi(tokenSplit[0])
	if err != nil {
		return m, err
	}

	return []dto.Bot{
		{
			ID:    botID,
			Token: botToken,
		},
	}, nil
}

func (d *dataProvider) GetChatsForBot(botID int) ([]dto.Chat, error) {
	chat := os.Getenv("SENDER_CHAT")
	c, _ := strconv.Atoi(chat)
	return []dto.Chat{{ID: int64(c)}}, nil
}

func (d *dataProvider) GetScheduledAd(year, month, day int) (dto.ScheduledAd, error) {
	return d.ad, nil
}

func (d *dataProvider) SetAd(ad dto.ScheduledAd) {
	d.ad = ad
}

func (d *dataProvider) LoadFile(path string) error {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	file := base64.StdEncoding.EncodeToString(fileBytes)
	d.ad.Message.FileBlob = &file

	return nil
}

func testMessage(sad dto.ScheduledAd, t *testing.T, fileToLoad string) {
	var err error

	fetcher := &dataProvider{}

	fetcher.SetAd(sad)
	if fileToLoad != "" {
		fetcher.LoadFile(fileToLoad)
	}

	bot, err := fetcher.GetBots()
	if err != nil {
		t.Error(err)
		return
	}

	chat := os.Getenv("SENDER_CHAT")
	c, _ := strconv.Atoi(chat)
	mailing, err := New(bot[0], fetcher, fetcher, dto.Chat{ID: int64(c)})
	if err != nil {
		t.Error(err)
		return
	}

	stop := make(chan struct{}, 100)

	// Try without reply markup
	mailing.Start(stop)

	// Try with reply markup
	sad.Message.ReplyMarkup = &tb.InlineKeyboardMarkup{
		InlineKeyboard: [][]tb.InlineButton{
			{{
				Text: "pohui",
				URL:  "https://google.com",
			}},
		},
	}
	fetcher.SetAd(sad)
	if fileToLoad != "" {
		fetcher.LoadFile(fileToLoad)
	}

	mailing, err = New(bot[0], fetcher, fetcher, dto.Chat{ID: int64(c)})
	if err != nil {
		t.Error(err)
		return
	}

	mailing.Start(stop)
}

func getDefaultAd() dto.ScheduledAd {
	now := time.Now()
	sad := dto.ScheduledAd{
		Year:    now.Year(),
		Month:   int(now.Month()),
		Day:     now.Day(),
		Message: dto.NewAd(),
	}
	sad.Message.Text = "<b>Just</b> <i>simple</i> <u>fucking</u> <s>text</s>"
	sad.Message.Test = false
	sad.Message.Interval = 1
	return sad
}

func TestText(t *testing.T) {
	sad := getDefaultAd()
	sad.Message.MessageType = "text"
	testMessage(sad, t, "")
}

func TestNewPhoto(t *testing.T) {
	sad := getDefaultAd()
	sad.Message.MessageType = "photo"
	testMessage(sad, t, "testdata/test_photo.jpg")
}

func TestNewVideo(t *testing.T) {
	sad := getDefaultAd()
	sad.Message.MessageType = "video"
	testMessage(sad, t, "testdata/test_video.mp4")
}

func TestNewAnimation(t *testing.T) {
	sad := getDefaultAd()
	sad.Message.MessageType = "animation"
	testMessage(sad, t, "testdata/test_animation.mp4")
}
