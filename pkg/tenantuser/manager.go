package tenantuser

import (
	"io"
	"log"
	"net/http"
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

type TenantUserKycReconciler struct {
	notifier.Client
	mgr *TenantUserKycManager
}

func (r *TenantUserKycReconciler) Reconcile(rkey interface{}) (*notifier.Result, error) {
	key, ok := rkey.(tenantuser.TenantUserKey)
	if !ok {
		log.Fatalln("Received key not of type Tenant user key in tenant user kyc manager: ", rkey)
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
	log.Printf("TenantUserMetadataReconciler: Entry fetched: %v\n", entry)
	baseUrl := r.mgr.baseUrl + "/" + entry.Key.Tenant + "/users/" + entry.Email + "/kyc"
	go func(url string) {
		if entry.Email != "" && entry.KYC == nil {
			log.Printf("KYC value not set for entry %s, fetching information from Tenant User KYC server\n", entry.Key.Username)
			// need to make an http call to fetch the value from db
			// getting client from package
			client := httpclient.GetClient()

			// making a request to the server
			// update the URL as per your requirement
			// for now not making it generic and will be making use of hardcoded values and appending at the end
			baseUrl := r.mgr.baseUrl + "/" + entry.Key.Tenant + "/users/" + entry.Email + "/kyc"
			req, err := http.NewRequest(r.mgr.httpMethod, baseUrl, nil)
			if err != nil {
				// something unexpected happened and we will retry after 5 seconds from the start
				log.Printf("Something unexpected went wrong while creating request,retrying again, error: %s\n", err)
				return
			}

			if r.mgr.apiKey != "" {
				req.Header.Set("apiKey", r.mgr.apiKey)
			}
			for {
				// Add a sleep interval to prevent 100% CPU usage
				time.Sleep(time.Duration(r.mgr.interval) * time.Second)
				log.Println("TenantUserKycReconciler: Sending request to fetch KYC information for tenant user:", entry.Key.Username)
				log.Println("baseUrl:", baseUrl)
				// Send the request
				resp, err := client.Do(req)
				if err != nil {
					// something unexpected happened and we will retry after 5 seconds from the start
					log.Printf("Something unexpected went wrong while sending request,retrying again, error: %s\n", err)
					continue
				}
				defer resp.Body.Close()

				// Read the response body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					// something unexpected happened and we will retry after 5 seconds from the start
					log.Printf("Something unexpected went wrong while reading response,retrying again, error: %s\n", err)
					continue
				}

				// Print the response for now, once APIs are provided this will be sent to specific package
				log.Println("Response status:", resp.Status)
				log.Println("Response body:", string(body))

				// convert the response into provider specific struct and fetch tenant user level kyc specific data
				// on basis of client name
				var kycStatus tenant.KYCStatus
				var kycErr error
				kycStatus, kycErr = r.mgr.provider.GetTenantUserKycStatus(body)
				if kycErr != nil {
					// something unexpected happened and we will retry after 5 seconds from the start
					log.Printf("Something unexpected went wrong while fetching kyb status from API,retrying again, error: %s\n", err)
					continue
				}

				// once value is fetched, need to set value in tenants collection
				// only setting value in case KYC was done successfully else keeping it as it is
				if kycStatus == tenant.KYCDone {
					// updating information in tenants collection
					update := &tenantuser.TenantUser{
						Key: entry.Key,
						KYC: &tenant.TenantKyc{
							Status: kycStatus,
						},
					}
					err := r.mgr.tenantUserTable.Update(update)
					if err != nil {
						// something unexpected happened and we will retry after 5 seconds from the start
						log.Printf("Something unexpected went wrong while updating kyb status in tenant config collection,retrying again, error: %s\n", err)
						continue
					}
				}

				// once value is set in db and config.stopOnceValueIsSet is true, send true to chan done
				if (kycStatus == tenant.KYCDone) && r.mgr.stopOnceSet {
					return
				}
				return
			}
		}
	}(baseUrl)

	return &notifier.Result{}, nil
}

type TenantUserKycManager struct {
	inframanager.ManagerImpl
	tenantUserTable *tenantuser.TenantUserTable
	interval        int
	stopOnceSet     bool
	enabled         bool
	baseUrl         string
	httpMethod      string
	apiKey          string
	provider        provider.Provider
}

func (m *TenantUserKycManager) Start() {
	r := &TenantUserKycReconciler{mgr: m}
	err := m.tenantUserTable.RegisterClient(TenantUserKycManagerClientName, r)
	if err != nil {
		log.Fatalln("failed to register tenant conf while starting TenantUserLevelKycManager", err)
	}
}

func CreateTenantUserKycManager(prvdr provider.Provider) *TenantUserKycManager {

	tenantUsrTbl, err := tenantuser.LocateTenantUserTable()
	if err != nil {
		log.Fatalln("unable to locate tenant config table:", err)
	}
	manager := &TenantUserKycManager{
		ManagerImpl: inframanager.ManagerImpl{
			InstanceKey: TenantUserKycManagerInstanceKey,
		},
		tenantUserTable: tenantUsrTbl,
		interval:        cfg.GetPollingTime(),
		enabled:         cfg.GetTenantUserMetadataEnabled(),
		baseUrl:         cfg.GetTenantUserMetadataEndpointDetails().BaseUrl,
		httpMethod:      cfg.GetTenantUserMetadataEndpointDetails().HttpMethod,
		apiKey:          cfg.GetTenantUserMetadataEndpointDetails().DefaultHeaders.Apikey,
		provider:        prvdr,
	}
	manager.InitImplWithTerminateHandling(manager)
	return manager
}
