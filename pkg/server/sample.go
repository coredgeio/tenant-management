package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/coredgeio/tenant-management/api/config"
)

type SampleApiServer struct {
	api.UnimplementedSampleApiServer
}

func NewSampleApiServer() *SampleApiServer {
	return &SampleApiServer{}
}

func (s *SampleApiServer) HelloWorld(ctx context.Context, req *api.HelloWorldReq) (*api.HelloWorldResp, error) {
	log.Println("got request", req)
	if req.Text == "error" {
		return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
	}
	resp := &api.HelloWorldResp{
		Project: req.Project,
		Offset:  req.Offset,
		Text:    req.Text,
	}
	return resp, nil
}
