package crud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type Language struct {
	ID           int    `json:"id"`
	UUID         string `json:"uuid"`
	ISOCode      string `json:"isoCode"`
	Title        string `json:"title"`
	MaleVoiceID  string `json:"maleVoiceId"`
	SearchConfig string `json:"searchConfig"`
}

type Op struct {
	Insert     string      `json:"insert"`
	Delete     int         `json:"delete,omitempty"`
	Retain     int         `json:"retain,omitempty"`
	Attributes interface{} `json:"attributes,omitempty"`
}

type Delta struct {
	Ops []*Op `json:"ops"`
}

type PathItem struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

type PublicUser struct {
	ID             int        `json:"id"`
	UUID           string     `json:"uuid"`
	Online         bool       `json:"online"`
	UserName       string     `json:"userName"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	Tutor          bool       `json:"tutor"`
	TutorOnline    bool       `json:"tutorOnline"`
	ImageMediaItem *MediaItem `json:"imageMediaItem"`
	Guest          bool       `json:"guest"`
}

type Textpad struct {
	Type           string      `json:"type"` // Enum: [ UNDEFINED, TEST, TEXTPAD, CANVAS, CANVAS_DOCUMENT, FOLDER, COLLECTION, LINK, WIKI, TASK_LIST, POST ]
	UUID           string      `json:"uuid"`
	Title          string      `json:"title"`
	Subtitle       string      `json:"subTitle"`
	User           *PublicUser `json:"user"`
	Owner          string      `json:"owner"`
	CreatedAt      int         `json:"createdAt,omitempty"`
	UpdatedAt      int         `json:"updatedAt,omitempty"`
	ImageMediaItem *MediaItem  `json:"imageMediaItem,omitempty"`
	// Tags []string `json:"tags"`
	Lang    *Language   `json:"language,omitempty"`
	Version int         `json:"version,omitempty"`
	Delta   *Delta      `json:"delta"`
	Path    []*PathItem `json:"path,omitempty"`
}

type ResponseTextpad struct {
	Data *Textpad `json:"data"`
}

// NewTextpad is a constructor for Textpad.
func NewTextpad(title, subtitle string, user *PrivateUser) *Textpad {
	return &Textpad{
		Type:     "TEXTPAD",
		UUID:     uuid.New().String(),
		Title:    title,
		Subtitle: subtitle,
		User:     user.ToPublic(),
		Owner:    user.UUID,
		Delta: &Delta{
			Ops: []*Op{
				{
					Insert: "\n",
				},
			},
		},
	}
}

// Create creates a textpad. X-Auth-Token is needed for this request
//
// path: /v1/textpads - /v1/library/{uuid}
//
// method: post
func (t *Textpad) Create(token string) (*ResponseTextpad, error) {
	client := &http.Client{}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(t)
	if err != nil {
		return nil, errors.Wrap(err, "encoding textpad data to json")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := acropliaAPIURL + fmt.Sprintf("/v1/library/%s", t.User.UUID)
	//path := acropliaAPIURL + "/v1/textpads"
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

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("request failed with status code %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	textpadResponse := &ResponseTextpad{}
	err = json.NewDecoder(resp.Body).Decode(textpadResponse)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return textpadResponse, nil
}
