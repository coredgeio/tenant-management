package provider

import "github.com/coredgeio/compass/controller/pkg/runtime/tenant"

type Provider interface {
	// functions to fetch different tenant level information

	// this function returns information on basis of whether KYB is done or not a tenant
	GetTenantKycStatus(body []byte) (tenant.KYCStatus, error)

	// this function returns information on basis of whether KYC is done or not for a tenant user
	GetTenantUserKycStatus(body []byte) (tenant.KYCStatus, error)

	// this function returns information on basis of whether tenant has configured a payment method or not for a tenant
	GetPaymentConfiguredStatus(body []byte) (bool, error)

	// this function returns information on basis of whether tenant type has been set or not for a tenant
	GetTenantType(body []byte) (tenant.TenantType, error)

	// this function returns the body in byte format which needs to be sent to provider to start metering
	BuildMeteringInfo(any) ([]byte, error)
}
