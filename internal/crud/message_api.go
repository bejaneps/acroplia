package crud

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var (
	// ErrMessageInternal is temporary error for chat message sending
	ErrMessageInternal = errors.New("request failed with status code 500 Internal Server Error, but message was sent")
)

type UserMessageEntity struct {
	UUID           string         `json:"uuid"`
	Type           string         `json:"type"` // Enum: UNDEFINED, USER_TEXT, USER_MEDIA, EVENT, DELETED
	Text           string         `json:"text"`
	Status         string         `json:"status"` // Enum: SENDING
	User           *PublicUser    `json:"user"`
	CreatedAt      int            `json:"createdAt,omitempty"`
	UpdatedAt      int            `json:"updatedAt,omitempty"`
	ReadMarks      map[string]int `json:"readMarks,omitempty"`
	Attachments    []struct{}     `json:"attachments"`
	ReplyToMessage string         `json:"replyToMessage,omitempty"`
}

type RespUserMessageEntity struct {
	Data *UserMessageEntity `json:"data"`
}

// NewMessage is a constructor for message to be sent in Acroplia
func NewMessage(text string, user *PrivateUser) *UserMessageEntity {
	return &UserMessageEntity{
		UUID:        uuid.New().String(),
		Type:        "USER_TEXT",
		Status:      "SENDING",
		User:        user.ToPublic(),
		Text:        text,
		Attachments: make([]struct{}, 0),
	}
}

// SendMessage sends message to specified chat uuid. X-Auth-Token is needed for this request.
//
// path: /v1/workspaces/{chat_uuid}/chat
//
// method: post
func (m *UserMessageEntity) SendMessage(chatUUID string, token string) (*RespUserMessageEntity, error) {
	client := &http.Client{}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(m)
	if err != nil {
		return nil, errors.Wrap(err, "encoding textpad data to json")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := acropliaAPIURL + fmt.Sprintf("/v1/workspaces/%s/chat", chatUUID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, buf)
	if err != nil {
		return nil, errors.Wrap(err, "creating new post request")
	}
	defer req.Body.Close()

	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"X-Auth-Token": []string{token},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "making post request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		return nil, ErrMessageInternal
	} else if resp.StatusCode != 500 {
		io.Copy(os.Stdout, resp.Body)
		return nil, errors.Errorf("request failed with status code %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	messageResponse := &RespUserMessageEntity{}
	err = json.NewDecoder(resp.Body).Decode(messageResponse)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return messageResponse, nil
}
