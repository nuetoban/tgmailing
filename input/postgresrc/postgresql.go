package postgresrc

import (
	"github.com/jmoiron/sqlx"

	"github.com/nuetoban/tgmailing/dto"
)

type Src struct {
	db     *sqlx.DB
	config Config
}

type Config struct {
	ChatsSelectQuery string
	BotsSelectQuery  string
	PostSelectQuery  string
}

func New(db *sqlx.DB, config Config) (*Src, error) {
	return &Src{
		db:     db,
		config: config,
	}, nil
}

func (s *Src) GetBots() ([]dto.Bot, error) {
	bots := []dto.Bot{}

	err := s.db.Select(&bots, s.config.BotsSelectQuery)

	return bots, err
}

func (s *Src) GetChatsForBot(botID int) ([]dto.Chat, error) {
	var chats []dto.Chat

	rows, err := s.db.NamedQuery(s.config.ChatsSelectQuery, map[string]interface{}{"bot_id": botID})
	for rows.Next() {
		var c dto.Chat
		err := rows.StructScan(&c)
		if err != nil {
			return chats, err
		}
		chats = append(chats, c)
	}

	return chats, err
}

func (s *Src) GetScheduledAd(y, m, d int) (dto.ScheduledAd, error) {
	var sad dto.ScheduledAd

	rows, err := s.db.NamedQuery(s.config.PostSelectQuery, map[string]interface{}{
		"year":  y,
		"month": m,
		"day":   d,
	})

	for rows.Next() {
		err := rows.StructScan(&sad)
		if err != nil {
			return sad, err
		}

		return sad, nil
	}

	return sad, err
}
