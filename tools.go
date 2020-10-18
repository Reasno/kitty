// +build tools

package tools

import (
	_ "github.com/gogo/protobuf/protoc-gen-gogo@21df5aa0e680850681b8643f0024f92d3b09930c"
	_ "github.com/gogo/protobuf/protoc-gen-gogofaster@21df5aa0e680850681b8643f0024f92d3b09930c"
	_ "github.com/gogo/protobuf/proto@21df5aa0e680850681b8643f0024f92d3b09930c"
	_ "github.com/kevinburke/go-bindata/go-bindata"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/metaverse/truss/truss"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/envoyproxy/protoc-gen-validate"
)
