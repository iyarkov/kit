package auth

import (
	"testing"
	"time"
)

func TestTokenReadWrite(t *testing.T) {
	token := Token{
		AccountId: 4,
		GroupId:   12,
		Role:      Operator,
		ExpiresAt: time.Now().Round(time.Second),
	}
	binaryToken := token.Write()

	var tokenCopy Token
	if err := tokenCopy.Read(binaryToken[:]); err != nil {
		t.Errorf("Unexpecte error %v", err)
	}

	if token != tokenCopy {
		t.Errorf("Tokens should match, token: %v, copy: %v", token, tokenCopy)
	}
}

func TestTokenReadInvalidBufferSize(t *testing.T) {
	binaryToken := [45]byte{}
	var token Token
	err := token.Read(binaryToken[:])
	if err != ErrorInvalidToken {
		t.Errorf("Expecting ErrorInvalidToken, got %v", err)
	}
}

func TestTokenReadInvalidFirstByte(t *testing.T) {
	binaryToken := [26]byte{}
	binaryToken[0] = 2
	var token Token
	err := token.Read(binaryToken[:])
	if err != ErrorInvalidToken {
		t.Errorf("Expecting ErrorInvalidToken, got %v", err)
	}
}

func TestTokenReadInvalidSecondByte(t *testing.T) {
	binaryToken := [26]byte{}
	binaryToken[0] = 2
	var token Token
	err := token.Read(binaryToken[:])
	if err != ErrorInvalidToken {
		t.Errorf("Expecting ErrorInvalidToken, got %v", err)
	}
}

func TestTokenIsExpiredTrue(t *testing.T) {
	token := Token{
		AccountId: 4,
		GroupId:   12,
		Role:      Operator,
		ExpiresAt: time.Now().Round(time.Second).Add(-time.Second),
	}

	if !token.IsExpired() {
		t.Errorf("Token must be expired")
	}
}

func TestTokenIsExpiredFalse(t *testing.T) {
	token := Token{
		AccountId: 4,
		GroupId:   12,
		Role:      Operator,
		ExpiresAt: time.Now().Round(time.Second).Add(2 * time.Second),
	}

	if token.IsExpired() {
		t.Errorf("Token must not be expired")
	}
}

func TestTokenReadWriteString(t *testing.T) {
	token := Token{
		AccountId: 4,
		GroupId:   12,
		Role:      Operator,
		ExpiresAt: time.Unix(1686841129, 0),
	}

	tokenString := token.WriteToString()

	tokenCopy, err := ReadFromString(tokenString)
	if err != nil {
		t.Errorf("Unexpecte error %v", err)
	}
	if token != *tokenCopy {
		t.Errorf("Tokens should match, token: %v, copy: %v", token, tokenCopy)
	}
}

func TestTokenReadInvalidString(t *testing.T) {
	_, err := ReadFromString("BDAAAAAAAAAABAAAAAwABgAAAAAAAGSLJyk=")
	if err != ErrorInvalidToken {
		t.Errorf("Expecting ErrorInvalidToken, got %v", err)
	}
}
