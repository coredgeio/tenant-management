module github.com/coredgeio/tenant-management

go 1.23.0

toolchain go1.23.8

require (
	github.com/coredgeio/compass v0.0.0-20250518125656-1aece1ccd516
	github.com/coredgeio/orbiter-auth v0.0.0-20250511064054-a7ecee4ecf78
	github.com/coredgeio/orbiter-baremetal-manager v0.0.0-20250326152721-03287621f7ae
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	go.mongodb.org/mongo-driver v1.12.1
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.5
	gopkg.in/yaml.v2 v2.4.0
)

replace cloud.google.com/go => cloud.google.com/go v0.100.2 // or a compatible version higher than v0.34.0

require (
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
)
