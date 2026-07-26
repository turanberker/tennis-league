package grpc

import (
	"context"
	"tennis-league/user-interface/grpc/pb/playerpb"
	"tennis-league/user-service/internal/service/player"
)

import (
	"errors"
	// Kendi 'user-interface' modülündeki üretilen pb paketini import ediyoruz
)

// PlayerGrpcServer gRPC isteklerini karşılayan yapıdır
type PlayerGrpcServer struct {
	playerpb.UnimplementedPlayerServiceServer                 // Go'nun ürettiği zorunlu taban struct
	playerUc                                  *player.Usecase // Mevcut UseCase katmanın
}

// NewPlayerGrpcServer yeni bir gRPC server nesnesi oluşturur
func NewPlayerGrpcServer(uc *player.Usecase) *PlayerGrpcServer {
	return &PlayerGrpcServer{
		playerUc: uc,
	}
}

// GetPlayer, proto dosyasında tanımladığımız metodun Go tarafındaki karşılığıdır
func (s *PlayerGrpcServer) GetPlayer(ctx context.Context, req *playerpb.GetPlayerRequest) (*playerpb.GetPlayerResponse, error) {
	playerID := req.GetPlayerId()
	if playerID == "" {
		return nil, errors.New("player_id bos olamaz")
	}

	// Senin mevcut player usecase katmanından oyuncuyu çekiyoruz
	p, err := s.playerUc.GetById(ctx, playerID) // NOT: UseCase'indeki metodun adı GetByID değilse kendi metoduna göre güncelle
	if err != nil {
		return nil, err
	}

	// Protobuf response tipinde veriyi dönüyoruz
	return &playerpb.GetPlayerResponse{
		Id:   p.ID,
		Name: p.Name + " " + p.Surname,
	}, nil
}
