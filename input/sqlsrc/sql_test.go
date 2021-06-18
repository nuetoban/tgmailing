package sqlsrc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuetoban/tgmailing/dto"
)

var db *sqlx.DB

func TestPostgreSQLConnectionParamsString(t *testing.T) {
	p := make(postgreSQLConnectionParams)

	p["dbname"] = "test_name"
	p["password"] = "test_password"
	p["user"] = "test_user"
	p["port"] = "5432"
	p["sslmode"] = "disable"

	expected := "dbname=test_name password=test_password port=5432 sslmode=disable user=test_user"

	if p.String() != expected {
		t.Errorf("wrong string: %s, expected: %s", expected, p)
	}
}

func TestGetBots(t *testing.T) {
	p, err := New(db, Config{BotsSelectQuery: "SELECT id, token FROM bots"})
	if err != nil {
		t.Errorf("New() returned error: %v", err)
		return
	}

	bots, err := p.GetBots()
	if err != nil {
		t.Errorf("GetBots() returned error: %v", err)
		return
	}

	expected := []dto.Bot{
		{ID: 123456, Token: "123456:foo"},
		{ID: 456789, Token: "456789:BAR"},
	}

	if !cmp.Equal(bots, expected) {
		t.Errorf("GetBots() returned wrong value: %v, expected: %v", bots, expected)
		return
	}
}

func TestGetChats(t *testing.T) {
	p, err := New(db, Config{ChatsSelectQuery: "SELECT id FROM chats"})
	if err != nil {
		t.Errorf("New() returned error: %v", err)
		return
	}

	chats, err := p.GetChatsForBot(0)
	if err != nil {
		t.Errorf("GetBots() returned error: %v", err)
		return
	}

	expected := []dto.Chat{
		{ID: -1},
		{ID: 0},
		{ID: 9_223_372_036_854_775_806},
	}

	if !cmp.Equal(chats, expected) {
		t.Errorf("GetChatsForBot() returned wrong value: %v, expected: %v", chats, expected)
		return
	}
}

func TestGetScheduledAd(t *testing.T) {
	p, err := New(db, Config{PostSelectQuery: "SELECT * FROM ads"})
	if err != nil {
		t.Errorf("New() returned error: %v", err)
		return
	}

	ad, err := p.GetScheduledAd(0, 0, 0)
	if err != nil {
		t.Errorf("GetBots() returned error: %v", err)
		return
	}

	expected := dto.ScheduledAd{}

	if !cmp.Equal(ad, expected) {
		t.Errorf("TestGetScheduledAd() returned wrong value: %v, expected: %v", ad, expected)
		return
	}
}

func TestMain(m *testing.M) {
	db = sqlx.MustConnect("sqlite3", ":memory:")

	db.MustExec("CREATE TABLE bots (id INTEGER, token TEXT)")
	db.MustExec("INSERT INTO bots (id, token) VALUES (123456, '123456:foo')")
	db.MustExec("INSERT INTO bots (id, token) VALUES (456789, '456789:BAR')")

	db.MustExec("CREATE TABLE chats (id BIGINT)")
	db.MustExec("INSERT INTO chats (id) VALUES (-1)")
	db.MustExec("INSERT INTO chats (id) VALUES (0)")
	db.MustExec("INSERT INTO chats (id) VALUES (9223372036854775806)")

	db.MustExec("CREATE TABLE ads (year INTEGER, month INTEGER, day INTEGER, message JSON)")
	db.MustExec(`INSERT INTO ads (year, month, day, message) VALUES
		(0, 0, 0, '{"message": {"test": false, "text": "Some <b>text</b>", "interval": 0.1, "message_type": "text"}}')`)

	m.Run()
}
