package tenantuserkyc

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coredgeio/compass/controller/pkg/runtime/tenant"
	pkgerrors "github.com/coredgeio/compass/pkg/errors"
	inframanager "github.com/coredgeio/compass/pkg/infra/manager"
	"github.com/coredgeio/compass/pkg/infra/notifier"
	"github.com/coredgeio/orbiter-auth/pkg/runtime/tenantuser"

	cfg "github.com/coredgeio/tenant-management/pkg/config"
	"github.com/coredgeio/tenant-management/pkg/httpclient"
	"github.com/coredgeio/tenant-management/pkg/provider"
)

const (
	// reconciler client for tenant table
	TenantUserKycManagerClientName = "TenantUserKycManager"
)

type TenantUserLevelKycReconciler struct {
	notifier.Client
	mgr *TenantUserLevelKycManager
}

func (r *TenantUserLevelKycReconciler) Reconcile(rkey interface{}) (*notifier.Result, error) {
	key, ok := rkey.(tenantuser.TenantUserKey)
	if !ok {
		log.Fatalln("Received key not of type Tenant user key in tenant user level kyc manager: ", rkey)
	}

	// check if manager lock is acquired
	if !r.mgr.IsOwnershipAcquired() {
		return &notifier.Result{}, nil
	}

	entry, err := r.mgr.tenantUserTable.Find(&key)
	if err != nil {
		if pkgerrors.IsNotFound(err) {
			log.Printf("entry not found for key: %s\n", key.Username)
			// nothing to do here, just exit
			return &notifier.Result{}, nil
		}
		// something unexpected happened and we will retry after 5 seconds from the start
		log.Println("Something unexpected went wrong,retrying again")
		return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
	}

	if entry.KYC == nil && r.mgr.enabled {
		log.Printf("KYC value not set for entry %s, fetching information from Tenant User Level KYC server\n", key.Username)
		for {
			// Add a sleep interval to prevent 100% CPU usage
			time.Sleep(time.Duration(r.mgr.interval) * time.Second)

			// need to make an http call to fetch the value from db
			// getting client from package
			client := httpclient.GetClient()

			// making a request to the server
			req, err := http.NewRequest(r.mgr.httpMethod, r.mgr.baseUrl, nil)
			if err != nil {
				// something unexpected happened and we will retry after 5 seconds from the start
				log.Printf("Something unexpected went wrong while creating request,retrying again, error: %s\n", err)
				return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
			}

			// Set headers
			req.Header.Set("Content-Type", r.mgr.contentTypeHeader)
			req.Header.Set("Authorization", r.mgr.authorizationHeader)
			req.Header.Set("apiKey", r.mgr.apiKey)

			// Send the request
			resp, err := client.Do(req)
			if err != nil {
				// something unexpected happened and we will retry after 5 seconds from the start
				log.Printf("Something unexpected went wrong while sending request,retrying again, error: %s\n", err)
				return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
			}
			defer resp.Body.Close()

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// something unexpected happened and we will retry after 5 seconds from the start
				log.Printf("Something unexpected went wrong while reading response,retrying again, error: %s\n", err)
				return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
			}

			// Print the response for now, once APIs are provided this will be sent to specific package
			log.Println("Response status:", resp.Status)
			log.Println("Response body:", string(body))

			// convert the response into provider specific struct and fetch kyb specific data
			// on basis of client name
			var prvder provider.Provider
			var kycStatus tenant.KYCStatus
			var kycErr error
			switch r.mgr.clientName {
			case "core42":
				prv := &provider.AiRev{
					ClientName: r.mgr.clientName,
				}
				prvder = prv
				kycStatus, kycErr = prvder.GetTenantLevelKycStatus(body)
			}

			if kycErr != nil {
				// something unexpected happened and we will retry after 5 seconds from the start
				log.Printf("Something unexpected went wrong while fetching kyb status from API,retrying again, error: %s\n", err)
				return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
			}

			// once value is fetched, need to set value in tenants collection
			// only setting value in case KYC was done successfully else keeping it as it is
			if kycStatus == tenant.KYCDone {
				r.mgr.mu.Lock()
				defer r.mgr.mu.Unlock()
				// updating information in tenants collection
				update := &tenantuser.TenantUser{
					Key: key,
					KYC: &tenant.TenantKyc{
						Status: kycStatus,
					},
				}
				err := r.mgr.tenantUserTable.Update(update)
				if err != nil {
					// something unexpected happened and we will retry after 5 seconds from the start
					log.Printf("Something unexpected went wrong while updating kyb status in tenant config collection,retrying again, error: %s\n", err)
					return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
				}
			}

			// once value is set in db and config.stopOnceValueIsSet is true, send true to chan done
			if (kycStatus == tenant.KYCDone) && r.mgr.stopOnceSet {
				break
			}
		}
	}

	return &notifier.Result{}, nil
}

type TenantUserLevelKycManager struct {
	inframanager.ManagerImpl
	tenantUserTable     *tenantuser.TenantUserTable
	interval            int
	stopOnceSet         bool
	enabled             bool
	baseUrl             string
	httpMethod          string
	contentTypeHeader   string
	authorizationHeader string
	clientName          string
	apiKey              string
	mu                  sync.Mutex
}

func (m *TenantUserLevelKycManager) Start() {
	r := &TenantUserLevelKycReconciler{mgr: m}
	err := m.tenantUserTable.RegisterClient(TenantUserKycManagerClientName, r)
	if err != nil {
		log.Fatalln("failed to register tenant conf while starting TenantUserLevelKycManager", err)
	}
}

func CreateTenantUserLevelKycManager() *TenantUserLevelKycManager {

	tenantUsrTbl, err := tenantuser.LocateTenantUserTable()
	if err != nil {
		log.Fatalln("unable to locate tenant config table:", err)
	}
	manager := &TenantUserLevelKycManager{
		ManagerImpl: inframanager.ManagerImpl{
			InstanceKey: TenantUserLevelKycManagerInstanceKey,
		},
		tenantUserTable:     tenantUsrTbl,
		interval:            cfg.GetTenantUserLevelKYCPollingTime(),
		enabled:             cfg.GetTenantUserLevelKYCEnabled(),
		baseUrl:             cfg.GetTenantUserLevelKYCBaseUrl(),
		httpMethod:          cfg.GetTenantUserLevelKYCHttpMethod(),
		contentTypeHeader:   cfg.GetTenantUserLevelKYCContentType(),
		authorizationHeader: cfg.GetTenantUserLevelKYCAuthorization(),
		apiKey:              cfg.GetTenantUserLevelKYCApiKey(),
		mu:                  sync.Mutex{},
	}
	manager.InitImplWithTerminateHandling(manager)
	return manager
}
