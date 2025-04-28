package paymentconfigured

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
	"github.com/coredgeio/compass/pkg/utils"

	cfg "github.com/coredgeio/tenant-management/pkg/config"
	"github.com/coredgeio/tenant-management/pkg/httpclient"
	"github.com/coredgeio/tenant-management/pkg/provider"
)

const (
	// reconciler client for tenant table
	PaymentConfiguredClientName = "PaymentConfiguredManager"
)

type PaymentConfiguredReconciler struct {
	notifier.Client
	mgr *PaymentConfiguredManager
}

func (r *PaymentConfiguredReconciler) Reconcile(rkey interface{}) (*notifier.Result, error) {
	key, ok := rkey.(tenant.TenantConfigKey)
	if !ok {
		log.Fatalln("Received key not of type domain config key in domain config manager: ", rkey)
	}

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

	// check if the entry already has kyc value set
	if entry.Kyc == nil && r.mgr.enabled {

		log.Printf("Payment configured value not set for entry %s, fetching information from Tenant config collection\n", key.Name)
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
			var paymentConfiguredStatus bool
			var paymentConfiguredStatusErr error
			switch r.mgr.clientName {
			case "core42":
				prv := &provider.AiRev{
					ClientName: r.mgr.clientName,
				}
				prvder = prv
				paymentConfiguredStatus, paymentConfiguredStatusErr = prvder.GetPaymentConfiguredStatus(body)
			}

			if paymentConfiguredStatusErr != nil {
				// something unexpected happened and we will retry after 5 seconds from the start
				log.Printf("Something unexpected went wrong while fetching payment configured status from API,retrying again, error: %s\n", err)
				return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
			}

			// once value is fetched, need to set value in tenants collection
			if paymentConfiguredStatus {
				r.mgr.mu.Lock()
				defer r.mgr.mu.Unlock()
				// updating information in tenants collection
				update := &tenant.TenantConfig{
					Key:               key,
					PaymentConfigured: utils.BoolP(paymentConfiguredStatus),
				}
				err := r.mgr.tenantConfTable.Update(update)
				if err != nil {
					// something unexpected happened and we will retry after 5 seconds from the start
					log.Printf("Something unexpected went wrong while updating payment configured status in tenant config collection, retrying again, error: %s\n", err)
					return &notifier.Result{NotifyAfter: 5 * time.Second}, nil
				}
			}

			// once value is set in db and config.stopOnceValueIsSet is true, send true to chan done
			if paymentConfiguredStatus && r.mgr.stopOnceSet {
				break
			}
		}

	}

	return &notifier.Result{}, nil
}

type PaymentConfiguredManager struct {
	inframanager.ManagerImpl
	tenantConfTable     *tenant.TenantConfigTable
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

func (m *PaymentConfiguredManager) Start() {
	r := &PaymentConfiguredReconciler{mgr: m}
	err := m.tenantConfTable.RegisterClient(PaymentConfiguredClientName, r)
	if err != nil {
		log.Fatalln("failed to register tenant conf while starting TenantManager", err)
	}
}

func CreatePaymentConfiguredManager() *PaymentConfiguredManager {

	tenantCfgTbl, err := tenant.LocateTenantConfigTable()
	if err != nil {
		log.Fatalln("unable to locate tenant config table:", err)
	}
	manager := &PaymentConfiguredManager{
		ManagerImpl: inframanager.ManagerImpl{
			InstanceKey: PaymentConfiguredManagerInstanceKey,
		},
		tenantConfTable:     tenantCfgTbl,
		interval:            cfg.GetPaymentMethodConfigurationPollingTime(),
		enabled:             cfg.GetPaymentMethodConfigurationEnabled(),
		baseUrl:             cfg.GetPaymentMethodConfigurationBaseUrl(),
		httpMethod:          cfg.GetPaymentMethodConfigurationHttpMethod(),
		contentTypeHeader:   cfg.GetPaymentMethodConfigurationContentType(),
		authorizationHeader: cfg.GetPaymentMethodConfigurationAuthorization(),
		apiKey:              cfg.GetPaymentMethodConfigurationApiKey(),
		mu:                  sync.Mutex{},
	}
	manager.InitImplWithTerminateHandling(manager)
	return manager
}
