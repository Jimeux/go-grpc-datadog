package rpc

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Jimeux/go-grpc-datadog/proto/go/pb/first/v1"
	"github.com/Jimeux/go-grpc-datadog/proto/go/pb/second/v1"
	"github.com/Jimeux/go-grpc-datadog/svc/first/internal/o11y"
)

type FirstService struct {
	first.UnimplementedFirstServiceServer

	secondSvcClient second.SecondServiceClient
}

func NewClientService(secondSvcClient second.SecondServiceClient) *FirstService {
	return &FirstService{secondSvcClient: secondSvcClient}
}

func (s *FirstService) Create(ctx context.Context, in *first.CreateRequest) (*first.CreateResponse, error) {
	o11y.Info(ctx, "start Create")
	if _, err := s.secondSvcClient.Create(ctx, &second.CreateRequest{Name: in.Name}); err != nil {
		o11y.Err(ctx, err, "Create failed for name="+in.Name)
		return nil, status.Error(codes.Unknown, "failed to create")
	}
	o11y.Info(ctx, "Create success")
	return &first.CreateResponse{Okay: true}, nil
}

func (s *FirstService) Fetch(ctx context.Context, in *first.FetchRequest) (*first.FetchResponse, error) {
	o11y.Info(ctx, "start Fetch")
	if _, err := s.secondSvcClient.GetByID(ctx, &second.GetByIDRequest{Id: in.Id}); err != nil {
		if status.Code(err) == codes.NotFound {
			return &first.FetchResponse{Okay: false}, nil
		}
		o11y.Err(ctx, err, "Fetch failed for id="+strconv.FormatInt(in.Id, 10))
		return nil, status.Error(codes.Unknown, "failed to fetch from ServerService")
	}
	o11y.Info(ctx, "Fetch success")
	return &first.FetchResponse{Okay: true}, nil
}
