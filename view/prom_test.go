package view

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nuetoban/tgmailing/dto"
)

type FakeStatsProvider struct{}

func (f *FakeStatsProvider) Statistics() dto.MultipleBotsStatistics {
	return dto.MultipleBotsStatistics{
		ID: 123456789123456789,
		Statistics: []dto.Statistics{
			{
				BotID:                  1,
				BotName:                "bot1",
				SuccessfulSendAttempts: 400,
				FailedSendAttempts:     666,
			},
			{
				BotID:                  2,
				BotName:                "bot2",
				SuccessfulSendAttempts: 777,
				FailedSendAttempts:     555,
			},
		},
		StartTime: time.Date(2020, 6, 25, 18, 0, 0, 0, time.UTC),
	}
}

type FakeHTTPResponceWriter struct {
	Headers http.Header
	Content strings.Builder
	Code    int
}

func (f *FakeHTTPResponceWriter) Header() http.Header {
	if f.Headers == nil {
		f.Headers = make(http.Header)
	}
	return f.Headers
}

func (f *FakeHTTPResponceWriter) Write(b []byte) (int, error) {
	return f.Content.Write(b)
}

func (f *FakeHTTPResponceWriter) WriteHeader(i int) {
	f.Code = i
}

func TestPromStatisticsHTTP(t *testing.T) {
	shouldContain := []string{
		"start_time{mailing_id=123456789123456789} 1593108000",
		"successful_send_attemts{bot_id=1,bot_name=bot1,mailing_id=123456789123456789} 400",
		"failed_send_attemts{bot_id=1,bot_name=bot1,mailing_id=123456789123456789} 666",
		"successful_send_attemts{bot_id=2,bot_name=bot2,mailing_id=123456789123456789} 777",
		"failed_send_attemts{bot_id=2,bot_name=bot2,mailing_id=123456789123456789} 555",
	}

	pp := &FakeStatsProvider{}
	fhttp := &FakeHTTPResponceWriter{}
	ph := NewPrometheus(pp)
	ph.StatisticsHTTP(fhttp, nil)

	out := fhttp.Content.String()

	for _, v := range shouldContain {
		if !strings.Contains(out, v) {
			t.Errorf("the output is not contain the string \"%s\". Output was:\n--- BEGIN\n%s\n--- END", v, out)
			return
		}
	}
}
