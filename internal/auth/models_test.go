package auth_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"

	"project/internal/auth"
)

func TestLoginInputDeserialize(t *testing.T) {
	data := `{"email": "test@email.com", "password": "password123"}`
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.LoginInput{}
	err := input.Deserialize(reader)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if input.Email != "test@email.com" || input.Password != "password123" {
		t.Errorf("Unexpected input: %+v", input)
	}
}

func TestLoginInputDeserializeError(t *testing.T) {
	data := `{"email": "test@email.com", "password": }` // malformed JSON
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.LoginInput{}
	err := input.Deserialize(reader)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestLoginInputDeserializeEmpty(t *testing.T) {
	data := `{"email": "", "password": ""}`
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.LoginInput{}
	err := input.Deserialize(reader)
	if err == nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if input.Email != "" || input.Password != "" {
		t.Errorf("Unexpected input: %+v", input)
	}
}

func TestRegisterInputDeserialize(t *testing.T) {
	data := `{"email": "test@email.com", "name": "Test User", "password": "password123"}`
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.RegisterInput{}
	err := input.Deserialize(reader)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if input.Email != "test@email.com" || input.Name != "Test User" || input.Password != "password123" {
		t.Errorf("Unexpected input: %+v", input)
	}
}

func TestRegisterInputDeserializeError(t *testing.T) {
	data := `{"email": "test@email.com", "name": "Test User", "password": }` // malformed JSON
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.RegisterInput{}
	err := input.Deserialize(reader)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestRegisterInputDeserializeEmpty(t *testing.T) {
	data := `{"email": "", "name": "", "password": ""}`
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.RegisterInput{}
	err := input.Deserialize(reader)
	if err == nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if input.Email != "" || input.Name != "" || input.Password != "" {
		t.Errorf("Unexpected input: %+v", input)
	}
}

func TestNewRegisterResp(t *testing.T) {
	input := &auth.RegisterInput{
		Email:    "test@email.com",
		Name:     "Test User",
		Password: "password123",
	}

	resp, err := auth.NewRegisterResp(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.UserID == uuid.Nil || resp.Registered.Before(time.Now().Add(-time.Minute)) || resp.UserStatus != "active" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestNewResgisterRespValues(t *testing.T) {
	input := &auth.RegisterInput{
		Email:    "",
		Name:     "",
		Password: "",
	}

	resp, err := auth.NewRegisterResp(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.UserID == uuid.Nil || resp.Registered.Before(time.Now().Add(-time.Minute)) || resp.UserStatus != "active" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestNewRegisterRespError(t *testing.T) {
	_, err := auth.NewRegisterResp(nil)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
