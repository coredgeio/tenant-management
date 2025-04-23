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

func (a *AiRev) GetKycStatus(body []byte) (tenant.KYCStatus, error) {

	return tenant.KYCDone, nil
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

func (a *AiRev) GetTenantType(body []byte) (string, error) {

	return "", nil
}
