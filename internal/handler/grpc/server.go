package grpc

import (
	"context"

	"github.com/FarhanDH/sidingin-rating-service/internal/platform/auth"
	"github.com/FarhanDH/sidingin-rating-service/internal/service"
	"github.com/FarhanDH/sidingin-rating-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedRatingServiceServer
	verifier *auth.JwtVerifier
	service  *service.RatingService
}

func NewGRPCServer(s *service.RatingService, v *auth.JwtVerifier) *GRPCServer {
	return &GRPCServer{
		verifier: v,
		service:  s,
	}
}

func (s *GRPCServer) GetBatchRatingStatus(ctx context.Context, req *pb.GetBatchRatingStatusRequest) (*pb.GetBatchRatingStatusResponse, error) {
	claims, err := s.verifier.VerifyIncoming(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}

	// get order IDs from request payload
	orderIDs := req.GetOrderIds()

	ratedMap, err := s.service.GetBatchRatingStatus(ctx, orderIDs, claims.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetBatchRatingStatusResponse{
		Statuses: ratedMap,
	}, nil
}
