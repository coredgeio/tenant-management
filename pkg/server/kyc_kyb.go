package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/coredgeio/tenant-management/api/config"
)

type TenantManagementServer struct {
	api.UnimplementedTenantManagementServer
}

func NewTenantManagementServer() *TenantManagementServer {
	return &TenantManagementServer{}
}

func (s *TenantManagementServer) GetKycStatus(ctx context.Context, req *api.KycStatusGetReq) (*api.KycStatusResp, error) {
	log.Printf("Received KYC status request for tenant: %s", req.GetName())

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant name is required")
	}

	// Sample logic — could fetch status from a DB or other service
	resp := &api.KycStatusResp{
		KycStatus: api.StatusInfo_KycDone, // returning an enum value
	}
	return resp, nil
}

func (s *TenantManagementServer) GetKybStatus(ctx context.Context, req *api.KybStatusGetReq) (*api.KybStatusResp, error) {
	log.Printf("Received KYB status request for tenant: %s", req.GetName())

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant name is required")
	}

	// Sample logic — could fetch status from a DB or other service
	resp := &api.KybStatusResp{
		KybStatus: api.StatusInfo_KycPending, // returning an enum value
	}
	return resp, nil
}
