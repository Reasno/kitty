# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.11 base image
FROM golang:1.15 as builder

# Set the Current Working Directory inside the container
WORKDIR /go/src/glab.tagtic.cn/ad_gains/kitty

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
ENV GONOPROXY *.tagtic.cn
ENV GOSUMDB sum.golang.google.cn
ENV GOPRIVATE glab.tagtic.cn/**

COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download dependencies
#RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/go-docker


######## Start a new stage from scratch #######
FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/' /etc/apk/repositories
RUN apk add ca-certificates tzdata

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY ./config ./config
COPY ./doc ./doc
COPY --from=builder /go/bin/go-docker .

EXPOSE 8080 9090

ENTRYPOINT ["./go-docker"]

CMD ["serve", "--config=./config/kitty.yaml"]
