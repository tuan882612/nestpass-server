package auth_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tuan882612/apiutils/securityutils"

	"project/internal/auth"
)

func Test_LoginInput_Deserialize_MissingParams(t *testing.T) {
	data := `{"email": "test@email.com", "password": "passwordpassword"}`
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.LoginInput{}
	err := input.Deserialize(reader)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func Test_LoginInput_Deserialize_Invalid(t *testing.T) {
	data := `{"email": "test@email.com", "password": }` // malformed JSON
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.LoginInput{}
	err := input.Deserialize(reader)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func Test_LoginInput_Deserialize_EmptyParams(t *testing.T) {
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

func Test_LoginInput_Deserialize_Nil(t *testing.T) {
	input := &auth.LoginInput{}
	if err := input.Deserialize(nil); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestRegisterInputDeserialize(t *testing.T) {
	data := `{"email": "test@email.com", "name": "Test User", "password": "passwordpassword123"}`
	reader := io.NopCloser(bytes.NewReader([]byte(data)))

	input := &auth.RegisterInput{}
	err := input.Deserialize(reader)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if input.Email != "test@email.com" || input.Name != "Test User" || input.Password != "passwordpassword123" {
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
		Password: "password123password123",
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
		Email:    "test@gmail.com",
		Name:     "test",
		Password: "passwordpassword123",
	}

	resp, err := auth.NewRegisterResp(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.Email != "test@gmail.com" || resp.Name != "test" {
		t.Errorf("Unexpected response: %+v", resp)
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

func TestNewRegisterRespEmtyInput(t *testing.T) {
	input := &auth.RegisterInput{
		Email:    "test@gmail.com",
		Name:     "",
		Password: "passwordpassword123",
	}

	_, err := auth.NewRegisterResp(input)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestNewRegisterRespNilInput(t *testing.T) {
	_, err := auth.NewRegisterResp(nil)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestNewRegisterRespProperHash(t *testing.T) {
	input := &auth.RegisterInput{
		Email:    "test@gmail.com",
		Name:     "testing",
		Password: "passwordpassword123",
	}

	resp, err := auth.NewRegisterResp(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := securityutils.ValidatePassword(resp.Password, input.Password); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
