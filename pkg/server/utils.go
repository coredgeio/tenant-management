package server

import (
	"github.com/coredgeio/compass/controller/pkg/runtime/tenant"
	api "github.com/coredgeio/tenant-management/api/config"
)

func mapKycStatus(status tenant.KYCStatus) api.KycStatus {
	var kyc api.KycStatus
	switch status {
	case tenant.KYCDone:
		kyc = api.KycStatus_Done
	case tenant.KYCFailed:
		kyc = api.KycStatus_Failed
	case tenant.KYCInProcess:
		kyc = api.KycStatus_InProcess
	case tenant.KYCPartial:
		kyc = api.KycStatus_Partial
	case tenant.ReKYCNeeded:
		kyc = api.KycStatus_ReKycNeeded
	case tenant.KYCPending:
		kyc = api.KycStatus_Pending
	default:
		kyc = api.KycStatus_Unspecified
	}
	return kyc
}
