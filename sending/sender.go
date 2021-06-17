package sending

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nuetoban/tgmailing/dto"

	tb "gopkg.in/tucnak/telebot.v2"
)

type DataFetcher interface {
	ChatsProvider
	AdProvider
}

type ChatsProvider interface {
	GetChatsForBot(botID int) ([]dto.Chat, error)
}

type AdProvider interface {
	GetScheduledAd(year, month, day int) (dto.ScheduledAd, error)
}

type Mailing struct {
	Statistics  dto.Statistics `json:"statistics"`
	ServiceChat dto.Chat

	// Hooks
	AfterStart  func(*Mailing) `json:"-"`
	AfterFinish func(*Mailing) `json:"-"`

	Bot     *tb.Bot         `json:"-"`       // telegram client
	Chats   []dto.Chat      `json:"chats"`   // list of chats to send post
	Retries []dto.Chat      `json:"retries"` // list of chats which did not recieve the post due to rate limits
	Post    dto.ScheduledAd `json:"post"`    // the post to send

	fileID string // file_id of uploaded file
}

// returns function which takes the chat and sends the message to telegram on call
func (m *Mailing) prepareSendableMessage() func(chat int64) error {
	var toSend tb.Sendable

	switch m.Post.Message.MessageType {
	case "text":

		return func(chat int64) error {
			var err error

			if m.Post.Message.ReplyMarkup != nil {
				_, err = m.Bot.Send(
					&tb.Chat{ID: chat},
					m.Post.Message.Text,
					tb.ModeHTML,
					&tb.ReplyMarkup{InlineKeyboard: m.Post.Message.ReplyMarkup.InlineKeyboard},
				)
			} else {
				_, err = m.Bot.Send(
					&tb.Chat{ID: chat},
					m.Post.Message.Text,
					tb.ModeHTML,
				)
			}

			return err
		}
	case "photo":
		toSend = &tb.Photo{File: tb.File{FileID: m.fileID}, Caption: m.Post.Message.Text}
	case "video":
		toSend = &tb.Video{File: tb.File{FileID: m.fileID}, Caption: m.Post.Message.Text}
	case "animation":
		toSend = &tb.Animation{
			File:     tb.File{FileID: m.fileID},
			Caption:  m.Post.Message.Text,
			FileName: "croco-ad.mp4",
		}
	default:
		return nil
	}

	return func(chat int64) error {
		var err error

		if m.Post.Message.ReplyMarkup != nil {
			_, err = m.Bot.Send(
				&tb.Chat{ID: chat},
				toSend,
				tb.ModeHTML,
				&tb.ReplyMarkup{InlineKeyboard: m.Post.Message.ReplyMarkup.InlineKeyboard},
			)
		} else {
			_, err = m.Bot.Send(
				&tb.Chat{ID: chat},
				toSend,
				tb.ModeHTML,
			)
		}

		return err
	}
}

// Start starts sending post to all chats
func (m *Mailing) Start(stop chan struct{}) {
	m.Statistics.StartTime = time.Now()

	log.Printf(
		"%20s 🤖 %26d 💚 %6d 💔 %6d 📚 --> Starting!\n",
		m.Bot.Me.Username, m.Statistics.SuccessfulSendAttempts, m.Statistics.FailedSendAttempts, len(m.Chats),
	)

	// Execute AfterStart hook
	if m.AfterStart != nil {
		m.AfterStart(m)
	}

	digits, _ := regexp.Compile("[0-9]+")

	for _, chat := range m.Chats {
		err := m.prepareSendableMessage()(chat.ID)

		if err != nil {
			m.Statistics.FailedSendAttempts++

			// If 429
			if strings.Contains(err.Error(), "Too Many Requests") {
				// Try to send later
				m.Retries = append(m.Retries, chat)
				var secondsToSleep int

				// Sleep the number of seconds which Telegram suggested to us
				if slp := digits.Find([]byte(err.Error())); slp != nil {
					secondsToSleep, _ = strconv.Atoi(string(slp))
				} else {
					// Sleep just 10
					secondsToSleep = 10
				}

				log.Printf(
					"%20s 🤖 %16d 💬 %6d 💚 %6d 💔 %6d 📚 --> 💔 Met 429, sleeping %d seconds\n",
					m.Bot.Me.Username, chat.ID, m.Statistics.SuccessfulSendAttempts, m.Statistics.FailedSendAttempts, len(m.Chats), secondsToSleep,
				)
				time.Sleep(time.Second * time.Duration(secondsToSleep))
			} else {
				log.Printf(
					"%20s 🤖 %16d 💬 %6d 💚 %6d 💔 %6d 📚 --> Failure 💔: %v\n",
					m.Bot.Me.Username, chat.ID, m.Statistics.SuccessfulSendAttempts, m.Statistics.FailedSendAttempts, len(m.Chats), err,
				)
			}
		} else {
			m.Statistics.SuccessfulSendAttempts++
			log.Printf(
				"%20s 🤖 %16d 💬 %6d 💚 %6d 💔 %6d 📚 --> Success 💚\n",
				m.Bot.Me.Username, chat.ID, m.Statistics.SuccessfulSendAttempts, m.Statistics.FailedSendAttempts, len(m.Chats),
			)
		}

		wait := time.Duration(m.Post.Message.Interval * 1000)
		time.Sleep(time.Millisecond * wait)
	}

	log.Printf(
		"%20s 🤖 %26d 💚 %6d 💔 %6d 📚 --> Starting retries (%d)\n",
		m.Bot.Me.Username, m.Statistics.SuccessfulSendAttempts, m.Statistics.FailedSendAttempts, len(m.Chats), len(m.Retries),
	)

	// Retry to send
	for _, chat := range m.Retries {
		err := m.prepareSendableMessage()(chat.ID)

		if err != nil {
			m.Statistics.FailedSendAttempts++

			log.Printf(
				"%20s 🤖 %16d 💬 %6d 💚 %6d 💔 %6d 📚 --> Failure 💔: %v\n",
				m.Bot.Me.Username, chat.ID, m.Statistics.SuccessfulSendAttempts, m.Statistics.FailedSendAttempts, len(m.Chats), err,
			)
		} else {
			m.Statistics.SuccessfulSendAttempts++

			log.Printf(
				"%20s 🤖 %16d 💬 %6d 💚 %6d 💔 %6d 📚 --> Success 💚\n",
				m.Bot.Me.Username, chat.ID, m.Statistics.SuccessfulSendAttempts, m.Statistics.FailedSendAttempts, len(m.Chats),
			)
		}

		wait := time.Duration(m.Post.Message.Interval * 1000)
		time.Sleep(time.Millisecond * wait)
	}

	log.Printf(
		"%20s 🤖 %26d 💚 %6d 💔 %6d 📚 --> Done 🎉\n",
		m.Bot.Me.Username, m.Statistics.SuccessfulSendAttempts, m.Statistics.FailedSendAttempts, len(m.Chats),
	)

	// Execute AfterFinish hook
	if m.AfterFinish != nil {
		m.AfterFinish(m)
	}

	stop <- struct{}{}
}

// SetAfterStartHook sets the function which will be executed after mailing start
func (m *Mailing) SetAfterStartHook(f func(*Mailing)) *Mailing {
	m.AfterStart = f
	return m
}

// SetAfterFinishHook sets the function which will be executed after mailing finish
func (m *Mailing) SetAfterFinishHook(f func(*Mailing)) *Mailing {
	m.AfterFinish = f
	return m
}

// New returns new instance of mailing
//
// One instance works for one bot.
// The function sets up telegram client, fetches the ad post and list of chats,
// and uploads media to telegram cloud if needed.
func New(bot dto.Bot, chatsFetcher ChatsProvider, adFetcher AdProvider, svcChat dto.Chat) (*Mailing, error) {
	var (
		err error
		m   = &Mailing{}
	)

	m.ServiceChat = svcChat

	// Setup telegram client
	log.Printf("%20d 🤖 🛠  Connecting to Telegram...\n", bot.ID)
	m.Bot, err = tb.NewBot(tb.Settings{Token: bot.Token})
	if err != nil {
		return m, err
	}

	m.Statistics.BotID = bot.ID
	m.Statistics.BotName = m.Bot.Me.Username

	username := m.Bot.Me.Username

	// Get ad post
	log.Printf("%20s 🤖 🛠  Fetching post...\n", username)
	now := time.Now()
	m.Post, err = adFetcher.GetScheduledAd(now.Year(), int(now.Month()), now.Day())
	if err != nil {
		return m, err
	}

	log.Printf("%20s 🤖 🛠  Fetching chats...\n", username)
	// If debug, send ad only to developers chat
	if m.Post.Message.Test {
		// Get developers
		m.Chats = []dto.Chat{{ID: svcChat.ID}}
	} else {
		// Get non-prem chats
		m.Chats, err = chatsFetcher.GetChatsForBot(bot.ID)
		if err != nil {
			return m, err
		}
	}
	log.Printf("%20s 🤖 🛠  Loaded %d chats\n", username, len(m.Chats))

	// If message is just text, no extra work required
	if m.Post.Message.MessageType == "text" {
		return m, nil
	}

	log.Printf("%20s 🤖 🛠  Uploading file\n", username)
	// Upload file to tg servers
	m.fileID, err = m.uploadFile(m.Post.Message.FileBlob)
	if err != nil {
		return m, err
	}
	time.Sleep(time.Second)

	return m, nil
}
