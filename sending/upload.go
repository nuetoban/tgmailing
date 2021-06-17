package sending

import (
	"bytes"
	"encoding/base64"
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

// uploads photo to telegram servers
func (m *Mailing) uploadPhoto(tgFile tb.File) (string, error) {
	finalFile := &tb.Photo{File: tgFile}
	message, err := m.Bot.Send(&tb.Chat{ID: m.ServiceChat.ID}, finalFile)
	if err != nil {
		return "", err
	}

	fid := message.Photo.MediaFile().FileID
	if fid == "" {
		return "", fmt.Errorf("file_id is empty after sending")
	}
	return fid, nil
}

// uploads video to telegram servers
func (m *Mailing) uploadVideo(tgFile tb.File) (string, error) {
	finalFile := &tb.Video{File: tgFile}
	message, err := m.Bot.Send(&tb.Chat{ID: m.ServiceChat.ID}, finalFile)
	if err != nil {
		return "", err
	}

	fid := message.Video.MediaFile().FileID
	if fid == "" {
		return "", fmt.Errorf("file_id is empty after sending")
	}
	return fid, nil
}

// uploads animation to telegram servers
func (m *Mailing) uploadAnimation(tgFile tb.File) (string, error) {
	finalFile := &tb.Animation{File: tgFile, FileName: "croco-ad.mp4"}
	message, err := m.Bot.Send(&tb.Chat{ID: m.ServiceChat.ID}, finalFile)
	if err != nil {
		return "", err
	}

	fid := message.Animation.MediaFile().FileID
	if fid == "" {
		return "", fmt.Errorf("file_id is empty after sending")
	}
	return fid, nil
}

// uploads file to telegram servers
func (m *Mailing) uploadFile(fileBlob *string) (string, error) {
	// Cannot be nil
	if fileBlob == nil {
		return "", fmt.Errorf("file_blob is nil")
	}

	// File in DB is base64-encoded, decode it
	decodedFile, err := base64.StdEncoding.DecodeString(*fileBlob)
	if err != nil {
		return "", err
	}

	// Understandable
	fileReader := bytes.NewReader(decodedFile)
	tgFile := tb.FromReader(fileReader)

	// Decide on the message type what to send to telegram
	switch m.Post.Message.MessageType {
	case "photo":
		return m.uploadPhoto(tgFile)
	case "video":
		return m.uploadVideo(tgFile)
	case "animation":
		return m.uploadAnimation(tgFile)
	default:
		return "", fmt.Errorf("not supported message_type")
	}
}
