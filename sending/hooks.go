package sending

import (
	"fmt"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// NotifyDevChatOnStart notifies about mailing start to dev chat
func NotifyDevChatOnStart(chatID int64) func(*SenderWithBots) {
	return func(s *SenderWithBots) {
		if len(s.Mailings) > 0 {
			m := s.Mailings[0]
			text := fmt.Sprintf("The mailing <code>%d</code> has been started!\n\n", s.ID)

			var (
				bots              []string
				chats             int
				maxChatsLen       int
				interval          time.Duration
				tgAverageRespTime = time.Millisecond * 58
			)

			for _, v := range s.Mailings {
				bots = append(bots, "@"+v.Bot.Me.Username)
				chats += len(v.Chats)
				if len(v.Chats) > maxChatsLen {
					maxChatsLen = len(v.Chats)
				}
				interval = time.Duration(int(v.Post.Message.Interval*1000)) * time.Millisecond
			}

			estimatedDuration := time.Duration(maxChatsLen) * (interval + tgAverageRespTime)

			text += fmt.Sprintf("<b>Bots</b>\n%s", strings.Join(bots, "\n"))
			text += fmt.Sprintf("\n\n<b>Loaded chats</b>\n%d", chats)
			text += fmt.Sprintf("\n\n<b>Estimated duration</b>\n%s", estimatedDuration)

			m.Bot.Send(
				&tb.Chat{ID: chatID},
				text,
				tb.ModeHTML,
			)
		}
	}
}

// NotifyDevChatOnFinishEach notifies about mailing finish to dev chat for each bot
func NotifyDevChatOnFinishEach(chatID int64) func(*Mailing) {
	return func(m *Mailing) {
		m.Bot.Send(
			&tb.Chat{ID: chatID},
			fmt.Sprintf(
				"The mailing for @%s has been finished!\n\n💚 %d\n💔 %d",
				m.Bot.Me.Username,
				m.Statistics.SuccessfulSendAttempts,
				m.Statistics.FailedSendAttempts,
			),
			tb.ModeHTML,
		)
	}
}

// DefaultAfterFinishEachHook notifies about mailing finish to dev chat for each bot
func NotifyDevChatOnFinish(chatID int64) func(*SenderWithBots) {
	return func(s *SenderWithBots) {
		if len(s.Mailings) > 0 {
			m := s.Mailings[0]

			var (
				totalSuccessful int
				totalFailed     int
			)

			for _, v := range s.Mailings {
				totalSuccessful += v.Statistics.SuccessfulSendAttempts
				totalFailed += v.Statistics.FailedSendAttempts
			}

			m.Bot.Send(
				&tb.Chat{ID: chatID},
				fmt.Sprintf(
					"The mailing for all bots has been finished!\n\n💚 %d\n💔 %d",
					totalSuccessful,
					totalFailed,
				),
				tb.ModeHTML,
			)
		}
	}
}
