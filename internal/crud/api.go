package crud

import (
	"bytes"
	"context"
	"net/http"
	"time"

	json "github.com/json-iterator/go"

	"github.com/pkg/errors"
)

type platform string

const (
	Undefined = platform("UNDEFINED")
	IOS       = platform("IOS")
	Android   = platform("ANDROID")

	acropliaAPIURL = "https://api-stage.acroplia.com/api"
)

var (
	// ErrInvalidCredentials is used when user submitted email, username, phone or password is invalid
	ErrInvalidCredentials = errors.New("invalid email, phone, username or password")
)

type JSONMediaDataSource struct {
	URL              string `json:"url"`
	FileSize         int    `json:"fileSize"`
	OriginalFileName string `json:"originalFileName"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
}

type JSONMediaDataJSONMediaDataSource struct {
	Type     string               `json:"type"` // Enum: IMAGE, AUDIO, VIDEO, DOCUMENT
	Metadata interface{}          `json:"metadata"`
	Source   *JSONMediaDataSource `json:"source"`
	Crop     string               `json:"crop"`
}

type MediaItem struct {
	ID   int                               `json:"id"`
	UUID string                            `json:"uuid"`
	Type string                            `json:"type"` // Enum: IMAGE, AUDIO, VIDEO, DOCUMENT
	Data *JSONMediaDataJSONMediaDataSource `json:"data"`
}

type CategoriesSettings struct {
	AccessChanging   string `json:"ACCESS_CHANGING"`   // Enum: ENABLED, DISABLED
	GroupAddContent  string `json:"GROUP_ADD_CONTENT"` // Enum: ENABLED, IMPORTANT_ONLY, DISABLED
	GroupMembership  string `json:"GROUP_MEMBERSHIP"`  // Enum: ENABLED, DISABLED
	GroupAssignments string `json:"GROUP_ASSIGNMENTS"` // Enum: ENABLED, DISABLED
	Mentions         string `json:"MENTIONS"`          // Enum: ENABLED, DISABLED
}

type UserNotificationsSettings struct {
	SMS         string              `json:"sms"`    // Enum: ENABLED, DISABLED
	Email       string              `json:"email"`  // Enum: ENABLED, DAYLY_DIGEST
	Stream      string              `json:"stream"` // Enum: ENABLED, DISABLED
	Categories  *CategoriesSettings `json:"categories"`
	MutedGroups []string            `json:"mutedGroups"`
}

type Role struct {
	ID   string `json:"id"`   // uuid
	Type string `json:"type"` // Enum: UNDEFINED, ANONYMOUS_USER, AUTHORIZED_USER, USER, GUEST, WORKSPACE_GUEST, WORKSPACE_MEMBER, WORKSPACE_ADMIN, NODE_OWNER, MODERATOR, PAID_MEMBER, BACK_OFFICER, NODE_BUYER, GXB_USER, COMMUNITY_GUEST, COMMUNITY_MEMBER, COMMUNITY_ADMIN, COMMUNITY_TASK_MANAGER
}

type PrivateUser struct {
	UUID                  string                     `json:"uuid"`
	ID                    int                        `json:"id"`
	Email                 string                     `json:"email"`
	Phone                 string                     `json:"phone"`
	RealName              string                     `json:"realName"`
	DisplayName           string                     `json:"displayName"`
	FirstName             string                     `json:"firstName"`
	LastName              string                     `json:"lastName"`
	ImageMediaItem        *MediaItem                 `json:"imageMediaItem"`
	UserName              string                     `json:"userName"`
	NotificationsSettings *UserNotificationsSettings `json:"notificationsSettings"`
	Roles                 []*Role                    `json:"roles"`
	BirthDate             int                        `json:"birthDate"`
	HasPassword           bool                       `json:"hasPassword"`
	IsOnline              bool                       `json:"isOnline"`
	IsGuest               bool                       `json:"isGuest"`
	IsTutor               bool                       `json:"isTutor"`
	IsTutorOnline         bool                       `json:"isTutorOnline"`
	IsActive              bool                       `json:"isActive"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Secret       string `json:"secret"`

	User *PrivateUser `json:"user"`
}

type ResponseAuthResponse struct {
	Data *AuthResponse `json:"data"`
}

// PushToken are tokens used to send push notification to target's device
type PushToken struct {
	FCMToken  string `json:"fcmToken"`
	APNSToken string `json:"apnsToken"`
}

// Device is a struct that holds info about user's device, such as IOS, Android and etc
type Device struct {
	UUID      string     `json:"uuid"`
	Platform  platform   `json:"platform"` // Enum: UNDEFINED, IOS, ANDROID
	Vendor    string     `json:"vendor"`
	Model     string     `json:"model"`
	OsVersion string     `json:"osVersion"`
	Locale    string     `json:"locale"`
	Tokens    *PushToken `json:"tokens"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// SignInEmailRequest is a struct that is ised for login by email requests
type SignInEmailRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Redirect string  `json:"redirect"`
	Device   *Device `json:"device"`
}

// NewSignInEmailRequest is a constructor for SignInEmailRequest type.
//
// It sets Device by default to null.
func NewSignInEmailRequest(email string, password string) *SignInEmailRequest {
	return &SignInEmailRequest{
		Email:    email,
		Password: password,
		Device:   nil,
	}
}

// SignInPhoneRequest is a struct that is used for login by phone requests
type SignInPhoneRequest struct {
	Phone    string  `json:"phone"`
	Password string  `json:"password"`
	Redirect string  `json:"redirect"`
	Device   *Device `json:"device"`
}

// NewSignInPhoneRequest is a constructor for SignInPhoneRequest type.
//
// It sets Device by default to null.
func NewSignInPhoneRequest(phone string, password string) *SignInPhoneRequest {
	return &SignInPhoneRequest{
		Phone:    phone,
		Password: password,
		Device:   nil,
	}
}

// SignInUsernameRequest is a struct that is ised for login by username requests
type SignInUsernameRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Redirect string  `json:"redirect"`
	Device   *Device `json:"device"`
}

// NewSignInUsernameRequest is a constructor for SignInUsernameRequest type.
//
// It sets Device by default to null.
func NewSignInUsernameRequest(username string, password string) *SignInUsernameRequest {
	return &SignInUsernameRequest{
		Username: username,
		Password: password,
		Device:   nil,
	}
}

// makeLoginRequests sends data in json format to specified login path
func makeLoginRequest(path string, data interface{}) (*ResponseAuthResponse, error) {
	client := &http.Client{}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return nil, errors.Wrap(err, "encoding user data to json")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", acropliaAPIURL+path, buf)
	if err != nil {
		return nil, errors.Wrap(err, "creating new post request")
	}
	defer req.Body.Close()

	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "making post request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidCredentials
	} else if resp.StatusCode != 200 {
		return nil, errors.Errorf("request failed with status code %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	acropliaResp := &ResponseAuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(acropliaResp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return acropliaResp, nil
}

// LoginByEmail makes a login request by email and password
//
// path: /v1/users/sessions
//
// method: post
func (s *SignInEmailRequest) LoginByEmail() (*ResponseAuthResponse, error) {
	resp, err := makeLoginRequest("/v1/users/sessions", s)
	if errors.Is(err, ErrInvalidCredentials) {
		return nil, ErrInvalidEmail
	}

	return resp, nil
}

// LoginByPhone makes a login request by phone and password
//
// path: /v1/users/sessions/phone
//
// method: post
func (s *SignInPhoneRequest) LoginByPhone() (*ResponseAuthResponse, error) {
	resp, err := makeLoginRequest("/v1/users/sessions/phone", s)
	if errors.Is(err, ErrInvalidCredentials) {
		return nil, ErrInvalidPhone
	}

	return resp, nil
}

// LoginByUsername makes a login request by username and password
//
// path: /v1/users/sessions/username
//
// method: post
func (s *SignInUsernameRequest) LoginByUsername() (*ResponseAuthResponse, error) {
	resp, err := makeLoginRequest("/v1/users/sessions/username", s)
	if errors.Is(err, ErrInvalidCredentials) {
		return nil, ErrInvalidUsername
	}

	return resp, nil
}
