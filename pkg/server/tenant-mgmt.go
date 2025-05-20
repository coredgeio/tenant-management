package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/coredgeio/compass/controller/pkg/runtime/tenant"
	"github.com/coredgeio/compass/pkg/auth"
	"github.com/coredgeio/compass/pkg/errors"
	"github.com/coredgeio/compass/pkg/utils"
	"github.com/coredgeio/orbiter-auth/pkg/runtime/tenantuser"
	api "github.com/coredgeio/tenant-management/api/config"
)

type TenantManagementServer struct {
	api.UnimplementedTenantMgmtApiServer
	tenantTable     *tenant.TenantConfigTable
	tenantUserTable *tenantuser.TenantUserTable
}

func NewTenantManagementServer() *TenantManagementServer {
	tbl, err := tenant.LocateTenantConfigTable()
	if err != nil {
		log.Fatalln("TenantManagementServer: unable to locate tenants collection: ", err)
	}
	tntUserTbl, err := tenantuser.LocateTenantUserTable()
	if err != nil {
		log.Fatalln("TenantManagementServer: unable to locate tenant user collection: ", err)
	}
	return &TenantManagementServer{
		tenantTable:     tbl,
		tenantUserTable: tntUserTbl,
	}
}

// fetches the tenant-kyc(kyb) status from tenants collection
func (s *TenantManagementServer) GetTenantKycStatus(ctx context.Context, req *api.GenericStatusReq) (*api.TenantKycStatusResp, error) {
	log.Printf("Received tenantKYC(KYB) status request for tenant")
	userInfo, found := auth.GetUserInfo(ctx)
	if !found {
		log.Println("user info not found")
		return nil, status.Error(codes.NotFound, "user info not found")
	}
	if userInfo == nil {
		log.Println("TenantManagementServer: GetTenantType: user info is nil")
		return nil, status.Error(codes.Internal, "user info is nil")
	}
	if userInfo.RealmName == "" {
		log.Println("TenantManagementServer: GetTenantType: realm name is missing")
		return nil, status.Error(codes.Internal, "tenant name is missing")
	}
	// logic to send tenant kyc status from tenants collection
	var tKycStatus bool
	tnt := s.tenantTable.Find(&tenant.TenantConfigKey{
		Name: userInfo.RealmName,
	})
	// return if tenant is not found
	if tnt == nil {
		return nil, status.Error(codes.NotFound, "tenant not found")
	}
	// check if tenant-kyc(KYB) status is set or not
	var kyc api.KycStatus
	if tnt.Kyc != nil {
		kyc = mapKycStatus(tnt.Kyc.Status)
		if tnt.Kyc.Status == tenant.KYCDone {
			tKycStatus = true
		}
	}
	resp := &api.TenantKycStatusResp{
		IsKycDone: tKycStatus,
		KycStatus: kyc,
	}
	return resp, nil
}

// fetches the tenant-user-kyc status from tenant-users collection
func (s *TenantManagementServer) GetTenantUserKycStatus(ctx context.Context, req *api.GenericStatusReq) (*api.TenantUserKycStatusResp, error) {
	log.Printf("Received KYC status request for tenant user")
	userInfo, found := auth.GetUserInfo(ctx)
	if !found {
		log.Println("user info not found")
		return nil, status.Error(codes.PermissionDenied, "Access denied, user info not found")
	}
	if userInfo == nil {
		log.Println("TenantManagementServer: GetTenantUserLevelKycStatus: user info is nil")
		return nil, status.Error(codes.Internal, "user info is nil")
	}
	if userInfo.UserName == "" || userInfo.RealmName == "" {
		log.Println("TenantManagementServer: GetTenantUserLevelKycStatus: user name or tenant name is missing")
		return nil, status.Error(codes.Internal, "user name or tenant name is missing")
	}
	// logic to send kyc status from tenant-users collection
	var kycStatus bool
	tntUser, err := s.tenantUserTable.Find(&tenantuser.TenantUserKey{
		Username: userInfo.UserName,
		Tenant:   userInfo.RealmName,
	})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Println("TenantManagementServer: GetTenantUserLevelKycStatus: tenant user not found")
			return nil, status.Error(codes.NotFound, "tenant user not found")
		}
		log.Println("TenantManagementServer: GetTenantUserLevelKycStatus: unable to find tenant user: ", err)
		return nil, status.Error(codes.Internal, "something unexpected happened while fetching tenant user")
	}
	// check if tenant user kyc status is present or not
	var kyc api.KycStatus
	if tntUser.KYC != nil {
		kyc = mapKycStatus(tntUser.KYC.Status)
		if tntUser.KYC.Status == tenant.KYCDone {
			kycStatus = true
		}
	}
	resp := &api.TenantUserKycStatusResp{
		Email:     tntUser.Email,
		Username:  tntUser.Key.Username,
		IsKycDone: kycStatus,
		KycStatus: kyc,
	}
	return resp, nil
}

// whether payment method is set or not for a tenant
func (s *TenantManagementServer) GetPaymentConfigStatus(ctx context.Context, req *api.GenericStatusReq) (*api.PaymentConfigStatusResp, error) {
	log.Printf("Received Payment Method Configured status request for tenant")
	userInfo, found := auth.GetUserInfo(ctx)
	if !found {
		log.Println("user info not found")
		return nil, status.Error(codes.PermissionDenied, "Access denied, user info not found")
	}
	if userInfo == nil {
		log.Println("TenantManagementServer: GetPaymentMethodConfigurationStatus: user info is nil")
		return nil, status.Error(codes.Internal, "user info is nil")
	}
	if userInfo.RealmName == "" {
		log.Println("TenantManagementServer: GetPaymentMethodConfigurationStatus: realm name is missing")
		return nil, status.Error(codes.Internal, "tenant name is missing")
	}
	// logic to send isPaymentMethodConf from tenants collection
	var isPaymentMethodConf bool
	tnt := s.tenantTable.Find(&tenant.TenantConfigKey{
		Name: userInfo.RealmName,
	})
	// return if tenant is not found
	if tnt == nil {
		return nil, status.Error(codes.NotFound, "tenant not found")
	}
	// check if tenant payment method is configured or not
	if tnt.PaymentConfigured != nil && utils.PBool(tnt.PaymentConfigured) {
		isPaymentMethodConf = true
	}
	resp := &api.PaymentConfigStatusResp{
		IsPayMethodSet: isPaymentMethodConf,
	}
	return resp, nil
}

// fetching tenant type, whether it is Individual or Organization
func (s *TenantManagementServer) GetTenantType(ctx context.Context, req *api.GenericStatusReq) (*api.TenantTypeResp, error) {
	log.Printf("Received Tenant Type request for tenant")
	userInfo, found := auth.GetUserInfo(ctx)
	if !found {
		log.Println("user info not found")
		return nil, status.Error(codes.PermissionDenied, "Access denied, user info not found")
	}
	if userInfo == nil {
		log.Println("TenantManagementServer: GetTenantType: user info is nil")
		return nil, status.Error(codes.Internal, "user info is nil")
	}
	if userInfo.RealmName == "" {
		log.Println("TenantManagementServer: GetTenantType: realm name is missing")
		return nil, status.Error(codes.Internal, "tenant name is missing")
	}
	// logic to send tenant type from tenants collection
	tnt := s.tenantTable.Find(&tenant.TenantConfigKey{
		Name: userInfo.RealmName,
	})
	// return if tenant is not found
	if tnt == nil {
		return nil, status.Error(codes.NotFound, "tenant not found")
	}
	// check if tenant type is present or not
	var tenantType api.TenantType
	if tnt.TenantType != nil {
		switch *tnt.TenantType {
		case tenant.Individual:
			tenantType = api.TenantType_Individual
		case tenant.Organisation:
			tenantType = api.TenantType_Organization
		default:
			tenantType = api.TenantType_Unknown
		}
	}
	resp := &api.TenantTypeResp{
		TenantType: tenantType,
	}
	return resp, nil
}
