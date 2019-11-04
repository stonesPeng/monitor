FROM golang:1.11.11 AS builder_base
ADD  ./src/go.mod /go/github.com/ZenLiu/GMonitor/go.mod
ADD  ./src/go.sum /go/github.com/ZenLiu/GMonitor/go.sum
WORKDIR /go/github.com/ZenLiu/GMonitor
ENV GOPROXY=https://athens.azurefd.net
RUN  GO111MODULE=on go mod download

FROM builder_base AS builder
ADD  ./src /go/github.com/ZenLiu/GMonitor
WORKDIR /go/github.com/ZenLiu/GMonitor
RUN  GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/main .
RUN ls -al /go/bin

FROM scratch
#FROM alpine
WORKDIR /app
COPY --from=builder  /go/bin/main ./main
#ADD ca-certificates.crt /etc/ssl/certs/
VOLUME ["/var/run/docker.sock","/app/config.yml","/usr/local/bin/docker"]
ENTRYPOINT  ["./main"]
