package sending

import (
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func TestRun(t *testing.T) {
	godotenv.Load(".env")
	godotenv.Load("../.env")

	dp := &dataProvider{}
	sad := getDefaultAd()
	sad.Message.MessageType = "text"
	dp.SetAd(sad)

	swb := NewSenderWithBots(dp, dp, dp)

	chat := os.Getenv("SENDER_CHAT")
	c, _ := strconv.Atoi(chat)

	// Set hooks
	swb.SetServiceChat(int64(c))
	swb.SetAfterStartHook(NotifyDevChatOnStart(int64(c)))
	swb.SetAfterFinishHook(NotifyDevChatOnFinish(int64(c)))
	swb.SetAfterFinishEachHook(NotifyDevChatOnFinishEach(int64(c)))

	err := swb.Run()
	if err != nil {
		t.Errorf("the SenderWithBots returned error: %v", err)
		return
	}

	stat := swb.Statistics()

	sa := stat.Statistics[0].SuccessfulSendAttempts
	if sa != 1 {
		t.Errorf("wrong SuccessfulSendAttempts count: %d", sa)
		return
	}

	fa := stat.Statistics[0].FailedSendAttempts
	if fa != 0 {
		t.Errorf("wrong FailedSendAttempts count: %d", fa)
		return
	}
}
