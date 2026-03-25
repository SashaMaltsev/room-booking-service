package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	httptransport "github.com/SashaMaltsev/room-booking-service/internal/transport/http"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Manager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

type Claims struct {
	UserID string            `json:"user_id"`
	Role   commondomain.Role `json:"role"`
	Exp    int64             `json:"exp"`
	Iat    int64             `json:"iat"`
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (m *Manager) Issue(userID string, role commondomain.Role) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", fmt.Errorf("%w: empty user id", ErrInvalidToken)
	}

	if !role.IsValid() {
		return "", fmt.Errorf("%w: invalid role", ErrInvalidToken)
	}

	now := m.now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		Exp:    now.Add(m.ttl).Unix(),
		Iat:    now.Unix(),
	}

	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	unsigned := encodeBase64URL(headerJSON) + "." + encodeBase64URL(claimsJSON)
	signature := m.sign(unsigned)

	return unsigned + "." + encodeBase64URL(signature), nil
}

func (m *Manager) Parse(token string) (httptransport.Principal, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return httptransport.Principal{}, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	signature, err := decodeBase64URL(parts[2])
	if err != nil {
		return httptransport.Principal{}, ErrInvalidToken
	}

	expected := m.sign(unsigned)
	if !hmac.Equal(signature, expected) {
		return httptransport.Principal{}, ErrInvalidToken
	}

	claimsBytes, err := decodeBase64URL(parts[1])
	if err != nil {
		return httptransport.Principal{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return httptransport.Principal{}, ErrInvalidToken
	}

	if claims.UserID == "" || !claims.Role.IsValid() {
		return httptransport.Principal{}, ErrInvalidToken
	}

	if claims.Exp <= m.now().Unix() {
		return httptransport.Principal{}, ErrTokenExpired
	}

	return httptransport.Principal{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}

func (m *Manager) sign(payload string) []byte {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(payload))
	return mac.Sum(nil)
}

func encodeBase64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func decodeBase64URL(value string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(value)
}
