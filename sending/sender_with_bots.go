package sending

import (
	"time"

	"github.com/nuetoban/tgmailing/dto"
)

// BotsProvider is created for fetching bots which will send mailing
type BotsProvider interface {
	GetBots() ([]dto.Bot, error)
}

// SenderWithBots sends mailing via multiple bots
type SenderWithBots struct {
	ID int64 `json:"id"`

	// Hooks
	AfterStart      func(*SenderWithBots) `json:"-"`
	AfterFinish     func(*SenderWithBots) `json:"-"`
	AfterStartEach  func(*Mailing)        `json:"-"`
	AfterFinishEach func(*Mailing)        `json:"-"`

	Mailings  []*Mailing `json:"mailings"`
	StartTime time.Time  `json:"start_time"`

	BotsProvider BotsProvider  `json:"-"`
	ChatsFetcher ChatsProvider `json:"-"`
	AdFetcher    AdProvider    `json:"-"`
	ServiceChat  dto.Chat      `json:"service_chat"`
}

// NewSenderWithBots returns new instance of mailing sender for multiple bots
func NewSenderWithBots(botsProvider BotsProvider, chatsFetcher ChatsProvider, adFetcher AdProvider) *SenderWithBots {
	return &SenderWithBots{
		ID:           time.Now().UnixNano(),
		BotsProvider: botsProvider,
		ChatsFetcher: chatsFetcher,
		AdFetcher:    adFetcher,
	}
}

// Statistics returns mailing metrics
func (s *SenderWithBots) Statistics() dto.MultipleBotsStatistics {
	stat := dto.MultipleBotsStatistics{}
	stat.StartTime = s.StartTime
	stat.ID = s.ID

	for _, m := range s.Mailings {
		stat.Statistics = append(stat.Statistics, m.Statistics)
	}

	return stat
}

// SetServiceChat sets the chat which will be used to upload files
func (s *SenderWithBots) SetServiceChat(chatID int64) *SenderWithBots {
	s.ServiceChat = dto.Chat{ID: chatID}
	return s
}

// SetAfterStartEachHook sets the function which will be executed after mailing start for each bot
func (s *SenderWithBots) SetAfterStartEachHook(f func(*Mailing)) *SenderWithBots {
	s.AfterStartEach = f
	return s
}

// SetAfterFinishEachHook sets the function which will be executed after mailing finish for each bot
func (s *SenderWithBots) SetAfterFinishEachHook(f func(*Mailing)) *SenderWithBots {
	s.AfterFinishEach = f
	return s
}

// SetAfterStartHook sets the function which will be executed after mailing start
func (s *SenderWithBots) SetAfterStartHook(f func(*SenderWithBots)) *SenderWithBots {
	s.AfterStart = f
	return s
}

// SetAfterFinishHook sets the function which will be executed after mailing finish
func (s *SenderWithBots) SetAfterFinishHook(f func(*SenderWithBots)) *SenderWithBots {
	s.AfterFinish = f
	return s
}

// Run starts the mailing. Stops when all bots complete its mailings.
func (s *SenderWithBots) Run() error {
	s.StartTime = time.Now()

	// Fetch bots which will send mailing
	bots, err := s.BotsProvider.GetBots()
	if err != nil {
		return err
	}

	// Create a chan to notify about mailings finish
	stop := make(chan struct{})

	// Create mailings for every bot
	for _, bot := range bots {
		m, err := New(bot, s.ChatsFetcher, s.AdFetcher, s.ServiceChat)
		if err != nil {
			return err
		}

		m.SetAfterStartHook(s.AfterStartEach)
		m.SetAfterFinishHook(s.AfterFinishEach)

		s.Mailings = append(s.Mailings, m)

		go m.Start(stop)
	}

	// Execute AfterStart hook if defined
	if s.AfterStart != nil {
		s.AfterStart(s)
	}

	// Wait for all bots to complete its mailings
	for i := 0; i < len(bots); i++ {
		<-stop
	}

	// Execute AfterFinish hook if defined
	if s.AfterFinish != nil {
		s.AfterFinish(s)
	}

	return nil
}
