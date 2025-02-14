package delivery

import (
	"context"

	"github.com/qosmioo/ulr-shortener/internal/storage"
	"github.com/qosmioo/ulr-shortener/internal/usecase"
	pb "github.com/qosmioo/ulr-shortener/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	usecase *usecase.URLShortenerService
	logger  *zap.Logger
	pb.UnimplementedGRPCHandlerServer
}

func NewGRPCServer(usecase *usecase.URLShortenerService, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{usecase: usecase, logger: logger}
}

func (s *GRPCServer) ShortenUrl(ctx context.Context, req *pb.ShortenUrlRequest) (*pb.ShortenUrlResponse, error) {
	shortURL, err := s.usecase.CreateShortURL(req.LongUrl)
	if err != nil {
		if err == storage.ErrURLExists {
			return nil, status.Error(codes.AlreadyExists, "URL already exists")
		}
		s.logger.Error("Failed to shorten URL", zap.Error(err))
		return nil, err
	}
	s.logger.Info("URL shortened", zap.String("longURL", req.LongUrl), zap.String("shortURL", shortURL))
	return &pb.ShortenUrlResponse{ShortUrl: shortURL}, nil
}

func (s *GRPCServer) GetUrl(ctx context.Context, req *pb.GetUrlRequest) (*pb.GetUrlResponse, error) {
	originalURL, err := s.usecase.GetOriginalURL(req.ShortUrl)
	if err != nil {
		s.logger.Error("Failed to get original URL", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Original URL retrieved", zap.String("shortURL", req.ShortUrl), zap.String("originalURL", originalURL))
	return &pb.GetUrlResponse{LongUrl: originalURL}, nil
}
