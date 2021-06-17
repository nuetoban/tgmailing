package view

import (
	"strings"
	"testing"
)

func TestJsonStatisticsHTTP(t *testing.T) {
	pp := &FakeStatsProvider{}
	fhttp := &FakeHTTPResponceWriter{}
	ph := NewJSON(pp)
	ph.StatisticsHTTP(fhttp, nil)

	expected := strings.ReplaceAll(`{"ID":123456789123456789,"Statistics":
[{"bot_id":1,"bot_name":"bot1","successful_send_attempts":400,
"failed_send_attemplts":666,"start_time":"0001-01-01T00:00:00Z"},
{"bot_id":2,"bot_name":"bot2","successful_send_attempts":777,
"failed_send_attemplts":555,"start_time":"0001-01-01T00:00:00Z"}],
"StartTime":"2020-06-25T18:00:00Z"}`, "\n", "")

	out := fhttp.Content.String()
	if out != expected {
		t.Errorf("test returned wrong json: %v, expected: %v", out, expected)
		return
	}

	if fhttp.Headers.Values("Content-Type")[0] != "application/json" {
		t.Errorf("test did not set Content-Type header")
		return
	}
}
