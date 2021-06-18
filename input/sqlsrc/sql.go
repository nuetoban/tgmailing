package sqlsrc

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/nuetoban/tgmailing/dto"
)

type Src struct {
	db     *sqlx.DB
	config Config
}

// Config contains queries to fetch data
type Config struct {
	ChatsSelectQuery string
	BotsSelectQuery  string
	PostSelectQuery  string
}

type postgreSQLConnectionParams map[string]string

func (p *postgreSQLConnectionParams) String() string {
	var params []string

	// Sort
	keys := make([]string, 0, len(*p))
	for k := range *p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		params = append(params, fmt.Sprintf("%s=%s", k, (*p)[k]))
	}

	return strings.Join(params, " ")
}

// NewPostgreSQL returns new PostgreSQL input source
//
// The function takes credentials from env:
//   - envPrefix+DB_USER
//   - envPrefix+DB_NAME
//   - envPrefix+DB_PASS
//   - envPrefix+DB_HOST
//   - envPrefix+DB_PORT
//
// Default envPrefix is "SENDER_"
func NewPostgreSQL(config Config, sslMode bool, envPrefix string) (*Src, error) {
	var (
		src    Src
		params = make(postgreSQLConnectionParams)
	)

	if envPrefix == "" {
		envPrefix = "SENDER_"
	}

	params["user"] = os.Getenv(envPrefix + "DB_USER")
	params["dbname"] = os.Getenv(envPrefix + "DB_NAME")
	params["password"] = os.Getenv(envPrefix + "DB_PASS")
	params["host"] = os.Getenv(envPrefix + "DB_HOST")
	params["port"] = os.Getenv(envPrefix + "DB_PORT")

	if params["port"] == "" {
		params["port"] = "5432"
	}

	if !sslMode {
		params["sslmode"] = "disable"
	}

	db, err := sqlx.Connect("postgres", params.String())
	if err != nil {
		return &src, err
	}

	return New(db, config)
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
