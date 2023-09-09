package database

import (
	"context"
	"os"
	"testing"

	"project/internal/config"
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
	pool, err := getPostgres(context.Background(), os.Getenv("PG_URL"), 2)
	if err != nil {
		t.Errorf("GetPostgres() error: %v", err)
	}

	if _, err := pool.Exec(context.Background(), `SELECT 1`); err != nil {
		t.Errorf("pool.Exec() error: %v", err)
	}
}

func Test_EmptyParamsPostgres(t *testing.T) {
	if _, err := getPostgres(context.Background(), "", 0); err == nil {
		t.Errorf("GetPostgres() error: %v", err)
	}

	if _, err := getRedis("", ""); err == nil {
		t.Errorf("GetRedis() error: %v", err)
	}
}

func Test_NewRedis(t *testing.T) {
	conn, err := getRedis(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PSW"))
	if err != nil {
		t.Errorf("GetRedis() error: %v", err)
	}

	if _, err := conn.Ping().Result(); err != nil {
		t.Errorf("conn.Ping() error: %v", err)
	}
}

func Test_EmptyParamsRedis(t *testing.T) {
	if _, err := getRedis("", ""); err == nil {
		t.Errorf("GetRedis() error: %v", err)
	}
}
