package proto

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Creates a new grpc connection to the server address.
func NewGRPCConn(address string) (*grpc.ClientConn, error) {
	log.Info().Msg("connecting to grpc server...")

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Err(err).Str("location", "NewGrpcConnection").Msg("failed to connect to grpc server")
		return nil, err
	}

	return conn, nil
}
