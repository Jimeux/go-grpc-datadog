package rpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Jimeux/go-grpc-datadog/proto/go/pb/second/v1"
	"github.com/Jimeux/go-grpc-datadog/svc/second/internal/db"
	"github.com/Jimeux/go-grpc-datadog/svc/second/internal/o11y"
)

type SecondService struct {
	second.UnimplementedSecondServiceServer

	dao db.DAO
}

func (s *SecondService) Create(ctx context.Context, in *second.CreateRequest) (*second.CreateResponse, error) {
	m, err := s.dao.Create(ctx, in.Name)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to create name=%s: %v", in.Name, err)
	}
	return &second.CreateResponse{Model: toProto(m)}, nil
}

func (s *SecondService) GetByID(ctx context.Context, in *second.GetByIDRequest) (*second.GetByIDResponse, error) {
	o11y.Info(ctx, "starting GetByID")
	m, err := s.dao.GetByID(ctx, in.Id)
	if err != nil {
		o11y.Err(ctx, err, "unexpected error")
		return nil, status.Errorf(codes.Unknown, "unexpected error for ID=%d: %v", in.Id, err)
	}
	if m == nil {
		o11y.Err(ctx, errors.New("not found err"), "GetByID not found")
		return nil, status.Errorf(codes.NotFound, "no model found for ID=%d", in.Id)
	}
	o11y.Info(ctx, "GetByID success")
	return &second.GetByIDResponse{Model: toProto(m)}, nil
}

func toProto(m *db.Model) *second.Model {
	return &second.Model{
		Id:        m.ID,
		Name:      m.Name,
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}
