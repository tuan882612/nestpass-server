package databases

import (
	"context"
	"os"
	"testing"

	"nestpass/internal/config"
)

func TestMain(m *testing.M) {
	err := config.LoadEnv("../../.env")
	if err != nil {
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func Test_NewPostgres(t *testing.T) {
	pool, err := NewPostgres(context.Background(), os.Getenv("PG_URL"))
	if err != nil {
		t.Errorf("NewPostgres() error: %v", err)
	}

	if _, err := pool.Exec(context.Background(), `SELECT 1`); err != nil {
		t.Errorf("pool.Exec() error: %v", err)
	}
}

func Test_EmptyParamsPostgres(t *testing.T) {
	if _, err := NewPostgres(context.Background(), ""); err == nil {
		t.Errorf("NewPostgres() error: %v", err)
	}

	if _, err := NewRedis("", ""); err == nil {
		t.Errorf("NewRedis() error: %v", err)
	}
}

func Test_NewRedis(t *testing.T) {
	conn, err := NewRedis(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PSW"))
	if err != nil {
		t.Errorf("NewRedis() error: %v", err)
	}

	if _, err := conn.Ping().Result(); err != nil {
		t.Errorf("conn.Ping() error: %v", err)
	}
}

func Test_EmptyParamsRedis(t *testing.T) {
	if _, err := NewRedis("", ""); err == nil {
		t.Errorf("NewRedis() error: %v", err)
	}
}
