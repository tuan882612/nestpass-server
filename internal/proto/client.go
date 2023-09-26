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
		log.Error().Str("location", "NewGrpcConnection").Msgf("failed to connect to grpc server: %v", err)
		return nil, err
	}

	return conn, nil
}
