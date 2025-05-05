package publishmeteringinfo

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/coredgeio/compass/controller/pkg/runtime/events"

	cfg "github.com/coredgeio/tenant-management/pkg/config"
	"github.com/coredgeio/tenant-management/pkg/httpclient"
	"github.com/coredgeio/tenant-management/pkg/provider"
)

type PublishMeteringInfoManager struct {
	events.EventsNotifier
	eventsTbl           *events.EventsTable
	baseUrl             string
	httpMethod          string
	contentTypeHeader   string
	authorizationHeader string
	apiKey              string
	provider            provider.Provider
}

func (m *PublishMeteringInfoManager) sendMeteringData(domain string, event *events.Event) {
	if event.Event.Resource == events.BMSResource {
		if event.Event.Data == nil {
			log.Println("No data to send for event", event.Event)
			return
		}
		resData := event.Event.Data.(*events.BMSResourceData)
		if resData.EventData == nil {
			log.Printf("EventData for BMS is nil; domain: %s | name: %s | id: %s", domain, resData.Name, resData.Id)
			return
		}
		var reqBody []byte
		switch resData.Event {
		case events.BMSAllocateEvent:
			// make an API call to send metering data
			reqBody, _ = m.provider.BuildMeteringInfo(events.BMSAllocateEvent)
		case events.BMSReleaseEvent:
			// make an API call to send metering data
			reqBody, _ = m.provider.BuildMeteringInfo(events.BMSReleaseEvent)
		default:
			log.Printf("Unknown event type for BMS metering data; domain: %s | name: %s | id: %s\n", domain, resData.Name, resData.Id)
			return
		}
		// need to make an http call to fetch the value from db
		// getting client from package
		client := httpclient.GetClient()
		// making a request to the server
		// update the URL as per your requirement
		// for now not making it generic and will be making use of hardcoded values and appending at the end
		req, err := http.NewRequest(m.httpMethod, m.baseUrl, bytes.NewReader(reqBody))
		if err != nil {
			// something unexpected happened
			log.Printf("Something unexpected went wrong while creating request, error: %s\n", err)
			return
		}
		// Set headers
		if m.contentTypeHeader != "" {
			req.Header.Set("Content-Type", m.contentTypeHeader)
		}
		if m.authorizationHeader != "" {
			req.Header.Set("Authorization", m.authorizationHeader)
		}
		if m.apiKey != "" {
			req.Header.Set("apiKey", m.apiKey)
		}
		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			// something unexpected happened
			log.Printf("Something unexpected went wrong while sending request, error: %s\n", err)
			return
		}
		defer resp.Body.Close()
		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// something unexpected happened
			log.Printf("Something unexpected went wrong while reading response, error: %s\n", err)
			return
		}
		// Print the response for now, once APIs are provided this will be sent to specific package
		log.Println("Response status:", resp.Status)
		log.Println("Response body:", string(body))

		// retry in case of error or not successful response
		// TODO[Akash] : This should be improved in future once we update manager implementation
		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to send metering data; domain: %s | name: %s | id: %s\n", domain, resData.Name, resData.Id)
			log.Printf("Retrying 3 more times...\n")
			for i := 0; i < 3; i++ {
				resp, err := client.Do(req)
				if err != nil {
					// something unexpected happened
					log.Printf("Something unexpected went wrong while sending request, error: %s\n", err)
					return
				}
				defer resp.Body.Close()
				// Read the response body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					// something unexpected happened
					log.Printf("Something unexpected went wrong while reading response, error: %s\n", err)
					return
				}
				// Print the response for now, once APIs are provided this will be sent to specific package
				log.Println("Response status:", resp.Status)
				log.Println("Response body:", string(body))
				if resp.StatusCode != http.StatusOK {
					i++
					continue
				} else {
					break
				}
			}
		}
		log.Printf("Successfully sent metering data; domain: %s | name: %s | id: %s\n", domain, resData.Name, resData.Id)
		// nothing to do, ignore for now
		return
	}
}

func (m *PublishMeteringInfoManager) OfflineUserNotify(domain string, key primitive.ObjectID, event *events.Event) {
	if event == nil {
		// nothing to do, ignore for now
		return
	}
	m.sendMeteringData(domain, event)
}

func CreatePublishMeteringInfoManager(prvdr provider.Provider) *PublishMeteringInfoManager {
	eventsTable, err := events.LocateEventsTable()
	if err != nil {
		log.Fatalln("Failed to located events table in publish metering info manager", err)
	}
	manager := &PublishMeteringInfoManager{
		eventsTbl:           eventsTable,
		baseUrl:             cfg.GetPublishMeteringInfoEndpointDetails().BaseUrl,
		httpMethod:          cfg.GetPublishMeteringInfoEndpointDetails().HttpMethod,
		contentTypeHeader:   cfg.GetPublishMeteringInfoEndpointDetails().DefaultHeaders.ContentType,
		authorizationHeader: cfg.GetPublishMeteringInfoEndpointDetails().DefaultHeaders.Authorization,
		apiKey:              cfg.GetPublishMeteringInfoEndpointDetails().DefaultHeaders.Apikey,
		provider:            prvdr,
	}
	eventsTable.RegisterClient(manager)
	return manager
}
