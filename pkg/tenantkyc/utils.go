package tenantkyc

import "github.com/coredgeio/compass/controller/pkg/runtime/tenant"

func AiRevMapKybStatusFromBoolToEnum(status bool) tenant.KYCStatus {
	if status {
		return tenant.KYCDone
	} else {
		return tenant.KYCPending
	}
}
