package tenant

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/coredgeio/compass/controller/pkg/runtime/tenant"
	pkgerrors "github.com/coredgeio/compass/pkg/errors"
	inframanager "github.com/coredgeio/compass/pkg/infra/manager"
	"github.com/coredgeio/compass/pkg/infra/notifier"
	"github.com/coredgeio/compass/pkg/utils"

	cfg "github.com/coredgeio/tenant-management/pkg/config"
	"github.com/coredgeio/tenant-management/pkg/httpclient"
	"github.com/coredgeio/tenant-management/pkg/provider"
)

const (
	// reconciler client for tenant table
	TenantMetadataManagerClientName = "TenantMetadataManager"
)

type TenantMetadataReconciler struct {
	notifier.Client
	mgr *TenantManager
}

/*
working of the reconciler function
1. Fetch the entry from db, if entry exists move forward else return.
2. if any one of kyc, paymentMethod or tenantType is not set, call the server to fetch these values.
3. start a loop and keep running it until kyc, paymentMethod and tenantType values are set and stopOnceSet is also true.
4. Inside loop, do not update all values with each request.
5. Update only those which are not present. Hence the if blocks inside the for loop.
*/
func (r *TenantMetadataReconciler) Reconcile(rkey interface{}) (*notifier.Result, error) {
	key, ok := rkey.(tenant.TenantConfigKey)
	if !ok {
		log.Fatalln("Received key not of type domain config key in domain config manager: ", rkey)
	}
	log.Printf("TenantMetadataReconciler: Received key: %s\n", key)

	// check if manager lock is acquired
	if !r.mgr.IsOwnershipAcquired() {
		return &notifier.Result{}, nil
	}

	entry, err := r.mgr.tenantConfTable.DBFind(&key)
	if err != nil {
		if pkgerrors.IsNotFound(err) {
			log.Printf("entry not found for key: %s\n", key.Name)
			// nothing to do here, just exit
			return &notifier.Result{}, nil
		}
		// something unexpected happened and we will retry after 5 seconds from the start
		log.Println("Something unexpected went wrong,retrying again")
		return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
	}
	log.Printf("TenantMetadataReconciler: Entry fetched: %v\n", entry)
	baseUrl := r.mgr.baseUrl + "/" + entry.Key.Name
	go func(baseUrl string) {
		if entry.Kyc == nil || entry.PaymentConfigured == nil || entry.TenantType == nil || (entry.Kyc != nil && entry.Kyc.Status != tenant.KYCDone) {
			log.Printf("KYC or Payment method or tenant type not set for entry %s, fetching information from external server\n", entry.Key.Name)
			// need to make an http call to fetch the value from db
			// getting client from package
			client := httpclient.GetClient()
			log.Printf("TenantMetadataReconciler: Client : %v\n", client)

			// making a request to the server
			// update the URL as per your requirement
			// for now not making it generic and will be making use of hardcoded values and appending at the end
			req, err := http.NewRequest(r.mgr.httpMethod, baseUrl, nil)
			if err != nil {
				// something unexpected happened and we will retry after 5 seconds from the start
				log.Printf("Something unexpected went wrong while creating request,retrying again, error: %s\n", err)
				return
			}
			log.Printf("TenantMetadataReconciler: Request created: %v\n", req)

			// Set headers
			if r.mgr.apiKey != "" {
				req.Header.Set("apiKey", r.mgr.apiKey)
			}

			for {
				log.Println("Request url: ", baseUrl)
				// Add a sleep interval to prevent 100% CPU usage
				time.Sleep(time.Duration(r.mgr.interval) * time.Second)

				// Send the request
				resp, err := client.Do(req)
				if err != nil {
					// something unexpected happened and we will retry after 5 seconds from the start
					log.Printf("Something unexpected went wrong while sending request,retrying again, error: %s\n", err)
					continue
				}
				defer resp.Body.Close()
				log.Printf("TenantMetadataReconciler: Response received: %v\n", resp)

				// Read the response body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					// something unexpected happened and we will retry after 5 seconds from the start
					log.Printf("Something unexpected went wrong while reading response,retrying again, error: %s\n", err)
					continue
				}
				log.Printf("TenantMetadataReconciler: Body received: %v\n", body)

				// Print the response for now, once APIs are provided this will be sent to specific package
				log.Println("TenantMetadataReconciler : Response status:", resp.Status)
				log.Println("TenantMetadataReconciler : Response body:", string(body))

				if resp.StatusCode != 200 {
					log.Printf("Invalid response: %s\n", err)
					continue
				}

				var tenantKycFound, tenantTypeFound, tenantPayMethodFound bool
				var currentKycStatus tenant.KYCStatus

				// check if the entry already has kyc value set
				if entry.Kyc == nil || (entry.Kyc != nil && entry.Kyc.Status != tenant.KYCDone) {
					// convert the response into provider specific struct and fetch tenant kyc specific data
					// on basis of client name
					var tenantKyc tenant.KYCStatus
					var kycErr error
					tenantKyc, kycErr = r.mgr.provider.GetTenantKycStatus(body)
					if kycErr != nil {
						// something unexpected happened and we will retry after 5 seconds from the start
						log.Printf("Something unexpected went wrong while fetching tenant kyc status from API,retrying again, error: %s\n", err)
						continue
					}

					// once value is fetched, need to set value in tenants collection
					// only setting value in case KYC was done successfully else keeping it as it is
					if tenantKyc == tenant.KYCDone || tenantKyc == tenant.KYCInProcess || tenantKyc == tenant.KYCFailed {
						if currentKycStatus == tenantKyc {
							log.Printf("KYC status for tenant %s is already set to %+v, skipping update\n", entry.Key.Name, tenantKyc)
							continue
						} else {
							currentKycStatus = tenantKyc
						}
						log.Printf("Updating KYC status for tenant: %s to : %+v", entry.Key.Name, tenantKyc)
						// updating information in tenants collection
						update := &tenant.TenantConfig{
							Key: entry.Key,
							Kyc: &tenant.TenantKyc{
								Status: tenantKyc,
							},
						}
						err := r.mgr.tenantConfTable.Update(update)
						if err != nil {
							// something unexpected happened and we will retry after 5 seconds from the start
							log.Printf("Something unexpected went wrong while updating tenant kyc status in tenant config collection,retrying again, error: %s\n", err)
							continue
						}
						if tenantKyc == tenant.KYCDone {
							tenantKycFound = true
						}
					}
				} else {
					tenantKycFound = true
				}
				// check if the entry already has tenant type value set
				if entry.TenantType == nil {
					// convert the response into provider specific struct and fetch tenant type specific data
					// on basis of client name
					var tenantType tenant.TenantType
					var tenantTypeErr error
					tenantType, tenantTypeErr = r.mgr.provider.GetTenantType(body)

					if tenantTypeErr != nil {
						// something unexpected happened and we will retry after 5 seconds from the start
						log.Printf("Something unexpected went wrong while fetching tenant kyc status from API,retrying again, error: %s\n", err)
						continue
					}
					log.Printf("Updating Tenant type for tenant: %s to : %+v", entry.Key.Name, tenantType)
					// once value is fetched, need to set value in tenants collection
					// only setting value in case TenantType was either Individual or organization
					// updating information in tenants collection
					update := &tenant.TenantConfig{
						Key:        entry.Key,
						TenantType: &tenantType,
					}
					err = r.mgr.tenantConfTable.Update(update)
					if err != nil {
						// something unexpected happened and we will retry after 5 seconds from the start
						log.Printf("Something unexpected went wrong while updating tenant kyc status in tenant config collection,retrying again, error: %s\n", err)
						continue
					}
					tenantTypeFound = true
				} else {
					tenantTypeFound = true
				}
				if entry.PaymentConfigured == nil {
					// convert the response into provider specific struct and fetch payment config data
					// on basis of client name
					var paymentConfiguredStatus bool
					var paymentConfiguredStatusErr error
					paymentConfiguredStatus, paymentConfiguredStatusErr = r.mgr.provider.GetPaymentConfiguredStatus(body)

					if paymentConfiguredStatusErr != nil {
						// something unexpected happened and we will retry after 5 seconds from the start
						log.Printf("Something unexpected went wrong while fetching payment configured status from API,retrying again, error: %s\n", err)
						continue
					}

					// once value is fetched, need to set value in tenants collection
					if paymentConfiguredStatus {
						log.Printf("Updating Payment Configured status for tenant: %s to : %+v", entry.Key.Name, paymentConfiguredStatus)
						// updating information in tenants collection
						update := &tenant.TenantConfig{
							Key:               entry.Key,
							PaymentConfigured: utils.BoolP(paymentConfiguredStatus),
						}
						err := r.mgr.tenantConfTable.Update(update)
						if err != nil {
							// something unexpected happened and we will retry after 5 seconds from the start
							log.Printf("Something unexpected went wrong while updating payment configured status in tenant config collection, retrying again, error: %s\n", err)
							continue
						}
						tenantPayMethodFound = true
					}
				} else {
					tenantPayMethodFound = true
				}
				// once value is set in db and config.stopOnceValueIsSet is true, send true to chan done
				if tenantKycFound && tenantPayMethodFound && tenantTypeFound && r.mgr.stopOnceSet {
					log.Println("Tenant : exiting goroutine : kyc done, payment done and tenant type done received and stoponceset : entry:", entry.Key.Name)
					return
				}

			}
		} else {
			return
		}
	}(baseUrl)

	return &notifier.Result{}, nil
}

type TenantManager struct {
	inframanager.ManagerImpl
	tenantConfTable *tenant.TenantConfigTable
	interval        int
	stopOnceSet     bool
	enabled         bool
	baseUrl         string
	httpMethod      string
	apiKey          string
	provider        provider.Provider
}

func (m *TenantManager) Start() {
	r := &TenantMetadataReconciler{mgr: m}
	err := m.tenantConfTable.RegisterClient(TenantMetadataManagerClientName, r)
	if err != nil {
		log.Fatalln("failed to register tenant conf while starting TenantMetadataManager", err)
	}
}

// setting up the manager
func CreateTenantMetadataManager(prvdr provider.Provider) *TenantManager {
	tenantCfgTbl, err := tenant.LocateTenantConfigTable()
	if err != nil {
		log.Fatalln("unable to locate tenant config table:", err)
	}
	manager := &TenantManager{
		ManagerImpl: inframanager.ManagerImpl{
			InstanceKey: TenantMetadataManagerInstanceKey,
		},
		tenantConfTable: tenantCfgTbl,
		interval:        cfg.GetPollingTime(),
		enabled:         cfg.GetTenantMetadataEnabled(),
		baseUrl:         cfg.GetTenantMetadataEndpointDetails().BaseUrl,
		httpMethod:      cfg.GetTenantMetadataEndpointDetails().HttpMethod,
		apiKey:          cfg.GetTenantMetadataEndpointDetails().DefaultHeaders.Apikey,
		provider:        prvdr,
	}
	manager.InitImplWithTerminateHandling(manager)
	return manager
}
