package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"forum/models"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

const sessionCookieName = "session"
const sessionSecretKey = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"

var UserSessions = make(map[string]string) //Key: sessionID, Value User

type SessionData struct {
	User      models.UserRegisteration `json:"user"`
	SessionID string                   `json:"sessionID"`
	ExpiresAt time.Time                `json:"expiresAt"`
}

func CreateSession(w http.ResponseWriter, r *http.Request, user models.UserRegisteration) error {
	sessionID := generateSessionID()
	expirationTime := time.Now().Add(24 * time.Hour) // Set session expiration to 24 hours

	sessionData := &SessionData{
		User:      user,
		SessionID: sessionID,
		ExpiresAt: expirationTime,
	}

	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	encodedSession := encodeSession(sessionJSON)

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    encodedSession,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true in production
		MaxAge:   int(24 * time.Hour / time.Second),
	}

	http.SetCookie(w, cookie)

	// Store the session ID in the UserSessions map
	UserSessions[user.Username] = sessionID

	return nil
}

func GetSessionUser(r *http.Request) (*models.UserRegisteration, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}

	decodedSession, err := decodeSession(cookie.Value)
	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	err = json.Unmarshal(decodedSession, &sessionData)
	if err != nil {
		return nil, err
	}

	// Check if the session has expired
	if time.Now().After(sessionData.ExpiresAt) {
		DestroySession(nil, r) // Destroy the session if it has expired
		return nil, errors.New("session has expired")
	}

	// Check if map has changed
	if sessionData.SessionID != UserSessions[sessionData.User.Username] {
		return &sessionData.User, errors.New("newer session found")
	}

	return &sessionData.User, nil
}

func DestroySession(w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true in production
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)

	// Remove the session from the UserSessions map
	if r != nil {
		cookie, err := r.Cookie(sessionCookieName)
		if err == nil {
			decodedSession, err := decodeSession(cookie.Value)
			if err == nil {
				var sessionData SessionData
				err = json.Unmarshal(decodedSession, &sessionData)
				if err == nil {
					delete(UserSessions, sessionData.User.Username)
				}
			}
		}
	}

	return nil
}

func encodeSession(data []byte) string {
	// Sign the data using HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(sessionSecretKey))
	mac.Write(data)
	expectedMAC := mac.Sum(nil)

	// Combine the data and the MAC
	combined := append(data, expectedMAC...)

	// Encode the combined data using base64
	encoded := base64.URLEncoding.EncodeToString(combined)

	return encoded
}

func decodeSession(encoded string) ([]byte, error) {
	// Decode the base64 encoded data
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// Validate the decoded data length
	if len(decoded) < sha256.Size {
		return nil, errors.New("invalid session data")
	}

	// Separate the data and the MAC
	dataLength := len(decoded) - sha256.Size
	data := decoded[:dataLength]
	providedMAC := decoded[dataLength:]

	// Calculate the expected MAC
	mac := hmac.New(sha256.New, []byte(sessionSecretKey))
	mac.Write(data)
	expectedMAC := mac.Sum(nil)

	// Compare the provided MAC with the expected MAC
	if !hmac.Equal(providedMAC, expectedMAC) {
		return nil, errors.New("invalid MAC")
	}

	return data, nil
}

func generateSessionID() string {
	// Generate a new UUID
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "default-session-id"
	}

	// Convert the UUID to a string
	return sessionID.String()
}
