package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/coredgeio/compass/controller/pkg/runtime/tenant"
	"github.com/coredgeio/compass/pkg/utils"
	api "github.com/coredgeio/tenant-management/api/config"
)

type TenantManagementServer struct {
	api.UnimplementedTenantMgmtApiServer
	tenantTable *tenant.TenantConfigTable
}

func NewTenantManagementServer() *TenantManagementServer {
	tbl, err := tenant.LocateTenantConfigTable()
	if err != nil {
		log.Fatalln("RegionApiServer: unable to locate region table: ", err)
	}
	return &TenantManagementServer{
		tenantTable: tbl,
	}
}

func (s *TenantManagementServer) GetTenantLevelKycStatus(ctx context.Context, req *api.GenericStatusReq) (*api.TenantLevelKycStatusResp, error) {
	log.Printf("Received KYB status request for tenant: %s", req.GetTenant())

	if req.GetTenant() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant name is required")
	}

	// logic to send kyb status from tenants collection
	var kybStatus bool
	tnt := s.tenantTable.Find(&tenant.TenantConfigKey{
		Name: req.GetTenant(),
	})
	// return if tenant is not found
	if tnt == nil {
		return nil, status.Error(codes.NotFound, "tenant not found")
	}
	// check if tenant KYB status is set or not
	if tnt.Kyc != nil || tnt.Kyc.Status == tenant.KYCDone {
		kybStatus = true
	}
	// Sample logic — could fetch status from a DB or other service
	resp := &api.TenantLevelKycStatusResp{
		Tenant:    req.GetTenant(),
		IsKybDone: kybStatus,
	}
	return resp, nil
}

// TODO[Akash]: this needs update due to recent changes in tenant user collection
func (s *TenantManagementServer) GetTenantUserLevelKycStatus(ctx context.Context, req *api.TenantUserLevelKycStatusReq) (*api.TenantUserLevelKycStatusResp, error) {
	log.Printf("Received KYB status request for tenant: %s, for user: %s", req.GetTenant(), req.GetEmail())

	if req.GetTenant() == "" || req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant name or email is missing")
	}

	var kycStatus bool

	// Sample logic — could fetch status from a DB or other service
	resp := &api.TenantUserLevelKycStatusResp{
		Tenant:    req.GetTenant(),
		Email:     req.GetEmail(),
		IsKycDone: kycStatus,
	}
	return resp, nil
}

func (s *TenantManagementServer) GetPaymentMethodConfigurationStatus(ctx context.Context, req *api.GenericStatusReq) (*api.PaymentMethodConfigurationStatusResp, error) {
	log.Printf("Received Payment Method Configured status request for tenant: %s", req.GetTenant())

	if req.GetTenant() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant name is missing")
	}

	var isPaymentMethodConf bool
	tnt := s.tenantTable.Find(&tenant.TenantConfigKey{
		Name: req.GetTenant(),
	})
	// return if tenant is not found
	if tnt == nil {
		return nil, status.Error(codes.NotFound, "tenant not found")
	}
	// check if tenant payment method is configured or not
	if tnt.PaymentConfigured != nil || utils.PBool(tnt.PaymentConfigured) {
		isPaymentMethodConf = true
	}

	// Sample logic — could fetch status from a DB or other service
	resp := &api.PaymentMethodConfigurationStatusResp{
		Tenant:             req.GetTenant(),
		IsPaymentMethodSet: isPaymentMethodConf,
	}
	return resp, nil
}

func (s *TenantManagementServer) GetTenantType(ctx context.Context, req *api.GenericStatusReq) (*api.TenantTypeResp, error) {
	log.Printf("Received Tenant Type request for tenant: %s", req.GetTenant())

	if req.GetTenant() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant name is missing")
	}

	var tenantType string
	tnt := s.tenantTable.Find(&tenant.TenantConfigKey{
		Name: req.GetTenant(),
	})
	// return if tenant is not found
	if tnt == nil {
		return nil, status.Error(codes.NotFound, "tenant not found")
	}
	// check if tenant type is present or not
	if tnt.PaymentConfigured != nil {
		if *tnt.TenantType == tenant.Individual {
			tenantType = "Individual"
		} else if *tnt.TenantType == tenant.Organisation {
			tenantType = "Organisation"
		} else {
			tenantType = "Unknown"
		}
	} else {
		tenantType = "Unknown"
	}

	// Sample logic — could fetch status from a DB or other service
	resp := &api.TenantTypeResp{
		Tenant:     req.GetTenant(),
		TenantType: tenantType,
	}
	return resp, nil
}
