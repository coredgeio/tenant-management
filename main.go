package main

import (
	"context"
	"flag"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/coredgeio/compass/controller/pkg/runtime/events"
	tenantruntime "github.com/coredgeio/compass/controller/pkg/runtime/tenant"
	"github.com/coredgeio/compass/pkg/auth"
	"github.com/coredgeio/compass/pkg/infra/configdb"
	"github.com/coredgeio/orbiter-auth/pkg/runtime/tenantuser"

	apiConfig "github.com/coredgeio/tenant-management/api/config"
	"github.com/coredgeio/tenant-management/api/config/swagger"
	"github.com/coredgeio/tenant-management/pkg/config"
	"github.com/coredgeio/tenant-management/pkg/provider"
	"github.com/coredgeio/tenant-management/pkg/publishmeteringinfo"
	"github.com/coredgeio/tenant-management/pkg/server"
	tenant "github.com/coredgeio/tenant-management/pkg/tenant"
	tuser "github.com/coredgeio/tenant-management/pkg/tenantuser"
)

const (
	// Internal GRPC port to host the grpc server
	GRPC_PORT = ":8090"

	// Port Over which registry API will be supported
	// for UI portal
	API_PORT = ":8080"
)

// parseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func parseFlags() (string, error) {
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "/opt/config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Return the configuration path
	return configPath, nil
}

func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(swagger.OpenAPI, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

func main() {
	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	err = config.ParseConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	err = configdb.InitializeDatabaseConnection(config.GetMongodbHost(),
		config.GetMongodbPort(), "compass-config")
	if err != nil {
		log.Println("Unable to initialize mongo database connection...")
		log.Println(err)
		log.Fatalln("Exiting...")
	}

	err = configdb.InitializeMetricsDatabaseConnection(config.GetMetricsdbHost(),
		config.GetMetricsdbPort(), "compass-metrics")
	if err != nil {
		log.Println("Unable to initialize metrics database connection...")
		log.Println(err)
		log.Fatalln("Exiting...")
	}

	// locating/initializing tenants collection
	_, err = tenantruntime.NewTenantConfigTable()
	if err != nil {
		log.Fatalln("unable to locate or create tenant config table")
	}

	// locating/initializing tenant-users collection
	_, err = tenantuser.LocateTenantUserTable()
	if err != nil {
		log.Fatalln("unable to locate or create tenant user table")
	}

	// locating/initializing events collection
	_, err = events.LocateEventsTable()
	if err != nil {
		log.Fatalln("Failed to locate events table", err)
	}

	// setting provider type
	var prvder provider.Provider
	switch config.GetClientName() {
	case "core42":
		prv := &provider.AiRev{
			ClientName: config.GetClientName(),
		}
		prvder = prv
	}

	// start the manager for tenant Level metatdata
	if config.GetTenantMetadataEnabled() {
		log.Println("Starting Tenant metadata manager...")
		// call tenant level manager which is working on notification from tenant collections
		// and updating the  tenant metadata at the tenant level collection
		go func() {
			_ = tenant.CreateTenantMetadataManager(prvder)
		}()
	}

	// start the manager for Tenant Type
	if config.GetTenantUserMetadataEnabled() {
		log.Println("Starting Tenant user level kyc manager...")
		// call tenant level manager which is working on notification from tenant collections
		// and updating the Tenant Type status at the tenant-user collection
		go func() {
			_ = tuser.CreateTenantUserKycManager(prvder)
		}()
	}

	// start the manager for publishing metering data
	if config.GetPublishMeteringInfoEnabled() {
		log.Println("Starting Publish Metering Info manager...")
		// call publish metering info manager which is working on notification from events collections
		// and sending data to third party to send to start metering
		go func() {
			_ = publishmeteringinfo.CreatePublishMeteringInfoManager(prvder)
		}()
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth.ProcessUserInfoInContext)))
	opts = append(opts, grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auth.ProcessUserInfoInContext)))
	grpcServer := grpc.NewServer(opts...)

	apiConfig.RegisterTenantMgmtApiServer(grpcServer, server.NewTenantManagementServer())

	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v for grpc server", err)
	}

	go func() {
		log.Println("serving grpc server...")
		log.Fatal(grpcServer.Serve(lis))
	}()

	// Create a client connection to just started grpc server
	conn, err := grpc.DialContext(
		context.Background(),
		GRPC_PORT,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("Failed to dial server:", err)
	}

	gwmux := gwruntime.NewServeMux(gwruntime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		// enable this section while using orbiter-auth module
		if key == auth.UserInfoHeader {
			return auth.UserInfoContext, true
		}
		return key, false
	}))

	// Register Tenant Managemnet API server
	err = apiConfig.RegisterTenantMgmtApiHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatal("Failed to register Tenant Management api handler with gateway:", err)
	}

	oa := getOpenAPIHandler()
	gwServer := &http.Server{
		Addr: API_PORT,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v") {
				gwmux.ServeHTTP(w, r)
				return
			}
			oa.ServeHTTP(w, r)
		}),
	}

	log.Println("starting api server")
	log.Println("Serving gRPC-Gateway on http://0.0.0.0" + API_PORT)
	log.Fatalln(gwServer.ListenAndServe())
}
