package email

import (
	"project/internal/config"
	"project/internal/proto"
	"project/internal/proto/pb/twofapb"
)

// Manager is used for sending different types of emails.
type Manager struct {
	Client twofapb.TwoFAServiceClient
}

// Used to initialize the email service client.
func NewManger(cfg *config.Configuration) (*Manager, error) {
	// initialize connection to grpc server
	conn, err := proto.NewGRPCConn(cfg.Server.Host + ":" + cfg.Server.GRPCPort)
	if err != nil {
		return nil, err
	}

	// load email service client
	return &Manager{Client: twofapb.NewTwoFAServiceClient(conn)}, nil
}
