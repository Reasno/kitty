module glab.tagtic.cn/ad_gains/kitty

go 1.14

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20200520232829-54ba9589114f
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.0 // indirect
	github.com/Reasno/trs v0.7.0
	github.com/antonmedv/expr v1.8.9
	github.com/aws/aws-sdk-go v1.29.5
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/gabriel-vasile/mimetype v1.1.1
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/go-kit/kit v0.10.0
	github.com/go-redis/redis/v8 v8.3.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/google/wire v0.4.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.0
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/knadh/koanf v0.14.0
	github.com/oklog/run v1.0.0
	github.com/opentracing-contrib/go-stdlib v1.0.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.3.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/cors v1.7.0
	github.com/rs/xid v1.2.1
	github.com/segmentio/kafka-go v0.4.8
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/spf13/cobra v1.1.0
	github.com/stretchr/testify v1.6.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	golang.org/x/tools v0.0.0-20201102212025-f46e4245211d // indirect
	google.golang.org/genproto v0.0.0-20201015140912-32ed001d685c
	google.golang.org/grpc v1.33.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gorm.io/driver/mysql v1.0.4-0.20201230025252-5b0f9700d29b
	gorm.io/driver/sqlite v1.1.1
	gorm.io/gorm v1.20.10-0.20210107034540-a5bfe2f39dab
)
