package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/coredgeio/compass/controller/pkg/runtime/tenant"
)

const (
	pending      = "PENDING"
	attempted    = "ATTEMPTED"
	rejected     = "REJECTED"
	failed       = "FAILED"
	individual   = "INDIVIDUAL"
	organization = "ORGANIZATION"
)

type AiRev struct {
	ClientName string
}

type genericResponse struct {
	Message string     `json:"message"`
	Data    tenantData `json:"data"`
}

type tenantData struct {
	ID       string `json:"id"`
	TenantId string `json:"tenantId"`
	// tenant kyc details
	KybDetails     kYBData     `json:"kyb"`
	PaymentDetails paymentData `json:"credit"`
	Tier           string      `json:"tier"`
	AccountType    string      `json:"accountType"`
}

type kYBData struct {
	Status string `json:"status"`
}

type tenantUserKycResponse struct {
	Message string            `json:"message"`
	Data    tenantUserKycData `json:"data"`
}

type tenantUserKycData struct {
	TenantId   string              `json:"tenantId"`
	Email      string              `json:"email"`
	KYCDetails tenantUserKYCStatus `json:"kyc"`
}

type tenantUserKYCStatus struct {
	Status string `json:"status"`
}

type paymentData struct {
	PaymentMethodAvailable bool `json:"paymentMethodAvailable"`
}

type metrics struct {
	CPU  int `json:"cpu"`
	RAM  int `json:"ram"`
	Disk int `json:"disk"`
}

type metadata struct {
	SystemID   string `json:"systemId"`
	Flavor     string `json:"flavor"`
	FlavorType string `json:"flavorType"`
}

type resourceEvent struct {
	Resource   string   `json:"resource"`
	Event      string   `json:"event"`
	TenantID   string   `json:"tenantId"`
	Region     string   `json:"region"`
	Domain     string   `json:"domain"`
	Project    string   `json:"project"`
	Name       string   `json:"name"`
	CreateTime int64    `json:"createTime"`
	CreatedBy  string   `json:"createdBy"`
	Metrics    metrics  `json:"metrics"`
	Metadata   metadata `json:"metadata"`
	Type       string   `json:"type"`
}

func (a *AiRev) GetTenantKycStatus(body []byte) (tenant.KYCStatus, error) {
	var resp genericResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Error while unmarshaling response for KYB in AiRev, error: %s\n", err)
		// using negative value to indicate that we have received an error response
		return -1, err
	}
	kybStatus := resp.Data.KybDetails.Status
	switch kybStatus {
	case pending:
		return tenant.KYCPending, nil
	case attempted:
		return tenant.KYCPending, nil
	case rejected:
		return tenant.KYCFailed, nil
	case failed:
		return tenant.KYCFailed, nil
	default:
		return tenant.KYCDone, nil
	}
}

func (a *AiRev) GetTenantUserKycStatus(body []byte) (tenant.KYCStatus, error) {

	var resp tenantUserKycResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Error while unmarshaling response for KYB in AiRev, error: %s\n", err)
		// using negative value to indicate that we have received an error response
		return -1, err
	}
	kycStatus := resp.Data.KYCDetails.Status
	switch kycStatus {
	case pending:
		return tenant.KYCPending, nil
	case attempted:
		return tenant.KYCPending, nil
	case rejected:
		return tenant.KYCFailed, nil
	case failed:
		return tenant.KYCFailed, nil
	default:
		return tenant.KYCDone, nil
	}
}

func (a *AiRev) GetPaymentConfiguredStatus(body []byte) (bool, error) {
	var resp genericResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Error while unmarshaling response for Payment Configured in AiRev, error: %s\n", err)
		return false, err
	}
	return resp.Data.PaymentDetails.PaymentMethodAvailable, nil
}

func (a *AiRev) GetTenantType(body []byte) (tenant.TenantType, error) {

	var resp genericResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Error while unmarshaling response for KYB in AiRev, error: %s\n", err)
		// using negative value to indicate that we have received an error response
		return -1, err
	}
	tenantType := resp.Data.AccountType
	switch tenantType {
	case individual:
		return tenant.Individual, nil
	default:
		return tenant.Organisation, nil
	}
}

func (a *AiRev) BuildMeteringInfo(data any) ([]byte, error) {

	event := resourceEvent{
		Resource:   "baremetal",
		Event:      "allocate",
		TenantID:   "tenant-abc",
		Region:     "us-central-1",
		Domain:     "prod",
		Project:    "project-x",
		Name:       "instance-1",
		CreateTime: time.Now().Unix(),
		CreatedBy:  "admin@example.com",
		Metrics: metrics{
			CPU:  4,
			RAM:  10240,
			Disk: 1000,
		},
		Metadata: metadata{
			SystemID:   "a8ftrh",
			Flavor:     "baremetal-large",
			FlavorType: "ComputeIntensive",
		},
		Type: "spot",
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	fmt.Println("JSON as byte array:", jsonBytes)
	fmt.Println("JSON string:", string(jsonBytes))
	return jsonBytes, nil
}
