package messaging

import (
	"context"
	"fmt"

	stream_chat "github.com/GetStream/stream-chat-go/v5"
)

// StreamService wraps the Stream Chat client with typed methods for BabyCare operations.
type StreamService struct {
	client *stream_chat.Client
}

// NewStreamService creates and verifies a Stream Chat client using the provided credentials.
func NewStreamService(apiKey, apiSecret string) (*StreamService, error) {
	client, err := stream_chat.NewClient(apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("stream: init client: %w", err)
	}
	return &StreamService{client: client}, nil
}

// UpsertUser creates or updates a user in Stream Chat.
// Must be called on user registration or login to keep Stream in sync.
func (s *StreamService) UpsertUser(userID, fullName string) error {
	_, err := s.client.UpsertUser(context.Background(), &stream_chat.User{
		ID:   userID,
		Name: fullName,
	})
	if err != nil {
		return fmt.Errorf("stream: upsert user %s: %w", userID, err)
	}
	return nil
}

// CreateChannel creates (or returns) a "messaging" channel for the given participants.
// channelID should be in the format "conversation-<parentID>-<babysitterID>".
func (s *StreamService) CreateChannel(channelID, parentID, babysitterID string) (*stream_chat.Channel, error) {
	resp, err := s.client.CreateChannel(
		context.Background(),
		"messaging",
		channelID,
		parentID,
		&stream_chat.ChannelRequest{
			Members: []string{parentID, babysitterID},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("stream: create channel %s: %w", channelID, err)
	}
	return resp.Channel, nil
}

// SendMessage sends a plain-text message to a channel on behalf of senderID.
func (s *StreamService) SendMessage(channelID, senderID, content string) error {
	ch := s.client.Channel("messaging", channelID)
	_, err := ch.SendMessage(
		context.Background(),
		&stream_chat.Message{
			Text: content,
			User: &stream_chat.User{ID: senderID},
		},
		senderID,
	)
	if err != nil {
		return fmt.Errorf("stream: send message to channel %s: %w", channelID, err)
	}
	return nil
}

// FreezeChannel sets the channel to frozen, preventing further messages.
func (s *StreamService) FreezeChannel(channelID string) error {
	ch := s.client.Channel("messaging", channelID)
	_, err := ch.Update(
		context.Background(),
		map[string]interface{}{"frozen": true},
		nil,
	)
	if err != nil {
		return fmt.Errorf("stream: freeze channel %s: %w", channelID, err)
	}
	return nil
}

// GetUnreadCount returns the number of messages sent to userID in channelID
// that the user has not yet read.
func (s *StreamService) GetUnreadCount(userID, channelID string) (int, error) {
	ch := s.client.Channel("messaging", channelID)

	resp, err := ch.Query(context.Background(), &stream_chat.QueryRequest{State: true})
	if err != nil {
		return 0, fmt.Errorf("stream: query channel %s: %w", channelID, err)
	}

	// Locate this user's last-read timestamp from the channel read state.
	var lastReadSet bool
	var lastReadUnix int64
	for _, r := range resp.Read {
		if r.User != nil && r.User.ID == userID {
			lastReadUnix = r.LastRead.Unix()
			lastReadSet = true
			break
		}
	}

	if !lastReadSet {
		// User has never read this channel — every message is unread.
		return len(resp.Messages), nil
	}

	// Count messages from other users that arrived after the last-read mark.
	count := 0
	for _, msg := range resp.Messages {
		if msg.User == nil || msg.User.ID == userID {
			continue
		}
		if msg.CreatedAt != nil && msg.CreatedAt.Unix() > lastReadUnix {
			count++
		}
	}
	return count, nil
}
