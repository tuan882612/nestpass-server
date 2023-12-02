package ping

import (
	"context"

	"github.com/rs/zerolog/log"

	"project/internal/config"
	"project/internal/proto"
	"project/internal/proto/pb/pingpb"
)

// PinManager used to ping the grpc parent server.
type PingManager struct {
	Client  pingpb.PingServiceClient
	PingReq *pingpb.PingData
}

// Used to initialize the ping service client.
func NewPingManager(cfg *config.Configuration) (*PingManager, error) {
	// initialize connection to grpc server
	conn, err := proto.NewGRPCConn(cfg.Server.Host + ":" + cfg.Server.GRPCPort)
	if err != nil {
		return nil, err
	}

	// initialize ping data
	pingData := &pingpb.PingData{
		Message: "auth server",
	}

	// load ping server client
	return &PingManager{
		Client:  pingpb.NewPingServiceClient(conn),
		PingReq: pingData,
	}, nil
}

func (m *PingManager) Ping() error {
	data, err := m.Client.Ping(context.Background(), m.PingReq)
	if err != nil {
		log.Error().Str("location", "Ping").Msgf("failed to ping server: %v", err)
		return err
	}

	log.Info().Str("location", "Ping").Msg(data.Message)
	return nil
}
