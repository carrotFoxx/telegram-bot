# STEP 1 build executable binary
FROM golang:1.15 AS builder
WORKDIR $GOPATH/src/iot-services/host-server
ADD ./ .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -v -o /go/bin/iotserver
RUN ls -a /go/bin

# STEP 2 build a small image
FROM scratch
COPY --from=builder /go/bin/iotserver /go/bin/iotserver
ENTRYPOINT ["/go/bin/iotserver"]