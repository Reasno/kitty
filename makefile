install:
	go mod tidy
	go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc \
        github.com/kevinburke/go-bindata/go-bindata \
        github.com/gogo/protobuf/protoc-gen-gogo \
        github.com/gogo/protobuf/protoc-gen-gogofaster \
        github.com/gogo/protobuf/proto \
        github.com/envoyproxy/protoc-gen-validate \

