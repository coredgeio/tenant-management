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

func (s *TenantManagementServer) GetTenantUserLevelKycStatus(ctx context.Context, req *api.TenantUserLevelKycGetReq) (*api.TenantUserLevelKycResp, error) {
	log.Printf("Received KYC status request for tenant: %s", req.GetName())

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant name is required")
	}

	// Sample logic — could fetch status from a DB or other service
	resp := &api.TenantUserLevelKycResp{
		TenantUserLevelKycStatus: api.StatusInfo_Done, // returning an enum value
	}
	return resp, nil
}

func (s *TenantManagementServer) GetTenantLevelKycStatus(ctx context.Context, req *api.TenantLevelKycGetReq) (*api.TenantLevelKycResp, error) {
	log.Printf("Received KYB status request for tenant: %s", req.GetName())

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant name is required")
	}

	// Sample logic — could fetch status from a DB or other service
	resp := &api.TenantLevelKycResp{
		TenantLevelKycStatus: api.StatusInfo_Pending, // returning an enum value
	}
	return resp, nil
}
