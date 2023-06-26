package auth

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	Anonymous Role = iota
	Admin
	Manager
	User
	Candidate
	Guest
	Operator
)

var ErrorInvalidToken = errors.New("invalid token")
var ErrorExpiredToken = errors.New("expired token")

type authTokenCtxKey struct{}

var bytOrder binary.ByteOrder = binary.BigEndian

type Role uint16

type Token struct {
	AccountId uint64
	GroupId   uint32
	Role      Role
	ExpiresAt time.Time
}

func (token *Token) IsInRole(role Role) bool {
	return token.Role == role
}

func (token *Token) IsAuthenticated() bool {
	return token.Role != Anonymous
}

func (token *Token) Write() [26]byte {
	var buffer [26]byte
	buffer[0] = 1
	bytOrder.PutUint64(buffer[2:10], token.AccountId)
	bytOrder.PutUint32(buffer[10:14], token.GroupId)
	bytOrder.PutUint16(buffer[14:18], uint16(token.Role))
	bytOrder.PutUint64(buffer[18:26], uint64(token.ExpiresAt.Unix()))
	return buffer
}

func (token *Token) Read(buffer []byte) error {
	if len(buffer) != 26 {
		return ErrorInvalidToken
	}
	if buffer[0] != 1 || buffer[1] != 0 {
		return ErrorInvalidToken
	}
	token.AccountId = bytOrder.Uint64(buffer[2:10])
	token.GroupId = bytOrder.Uint32(buffer[10:14])
	token.Role = Role(bytOrder.Uint16(buffer[14:18]))
	token.ExpiresAt = time.Unix(int64(bytOrder.Uint64(buffer[18:26])), 0)
	return nil
}

func (token *Token) IsExpired() bool {
	if !token.IsAuthenticated() {
		return false
	}
	return time.Now().After(token.ExpiresAt)
}

func (token *Token) WriteToString() string {
	bytes := token.Write()
	return base64.StdEncoding.EncodeToString(bytes[:])
}

func ReadFromString(encoded string) (*Token, error) {
	buffer, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode string %w", err)
	}
	token := Token{}
	return &token, token.Read(buffer)
}

func WithToken(ctx context.Context, tokenRef *Token) context.Context {
	return context.WithValue(ctx, &authTokenCtxKey{}, tokenRef)
}

func WithStringToken(ctx context.Context, token string) (context.Context, error) {
	tokenRef, err := ReadFromString(token)
	if err != nil {
		return nil, err
	}
	if tokenRef.IsExpired() {
		return nil, ErrorExpiredToken
	}
	return context.WithValue(ctx, &authTokenCtxKey{}, tokenRef), err
}

func AuthToken(ctx context.Context) Token {
	if token, ok := ctx.Value(&authTokenCtxKey{}).(*Token); ok {
		return *token
	}
	return Token{}
}
