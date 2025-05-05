package provider

import (
	"encoding/json"
	"log"

	"github.com/coredgeio/compass/controller/pkg/runtime/tenant"
)

type AiRev struct {
	ClientName string
}

type GenericResponse struct {
	Message string     `json:"message"`
	Data    TenantData `json:"data"`
}

type TenantData struct {
	ID             string      `json:"id"`
	TenantId       string      `json:"tenantId"`
	KybDetails     KYBData     `json:"kyb"`
	PaymentDetails PaymentData `json:"credit"`
	Tier           string      `json:"tier"`
	AccountType    string      `json:"accountType"`
}

type KYBData struct {
	Status string `json:"status"`
}

type TenantUserKycResponse struct {
	Message string            `json:"message"`
	Data    TenantUserKycData `json:"data"`
}

type TenantUserKycData struct {
	TenantId   string              `json:"tenantId"`
	Email      string              `json:"email"`
	KYCDetails TenantUserKYCStatus `json:"kyc"`
}

type TenantUserKYCStatus struct {
	Status string `json:"status"`
}

type PaymentData struct {
	PaymentMethodAvailable bool `json:"paymentMethodAvailable"`
}

func (a *AiRev) GetTenantLevelKycStatus(body []byte) (tenant.KYCStatus, error) {
	var resp GenericResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Error while unmarshaling response for KYB in AiRev, error: %s\n", err)
		return 100, err
	}
	kybStatus := resp.Data.KybDetails.Status
	switch kybStatus {
	case "PENDING":
		return tenant.KYCPending, nil
	case "ATTEMPTED":
		return tenant.KYCPending, nil
	case "REJECTED":
		return tenant.KYCFailed, nil
	case "FAILED":
		return tenant.KYCFailed, nil
	default:
		return tenant.KYCDone, nil
	}
}

func (a *AiRev) GetTenantUserLevelKycStatus(body []byte) (tenant.KYCStatus, error) {

	var resp TenantUserKycResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Error while unmarshaling response for KYB in AiRev, error: %s\n", err)
		return 100, err
	}
	kycStatus := resp.Data.KYCDetails.Status
	switch kycStatus {
	case "PENDING":
		return tenant.KYCPending, nil
	case "ATTEMPTED":
		return tenant.KYCPending, nil
	case "REJECTED":
		return tenant.KYCFailed, nil
	case "FAILED":
		return tenant.KYCFailed, nil
	default:
		return tenant.KYCDone, nil
	}
}

func (a *AiRev) GetPaymentConfiguredStatus(body []byte) (bool, error) {
	var resp GenericResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Error while unmarshaling response for Payment Configured in AiRev, error: %s\n", err)
		return false, err
	}
	return resp.Data.PaymentDetails.PaymentMethodAvailable, nil
}

func (a *AiRev) GetTenantType(body []byte) (tenant.TenantType, error) {

	var resp GenericResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Error while unmarshaling response for KYB in AiRev, error: %s\n", err)
		return 100, err
	}
	tenantType := resp.Data.AccountType
	switch tenantType {
	case "INDIVIDUAL":
		return tenant.Individual, nil
	default:
		return tenant.Organisation, nil
	}
}
