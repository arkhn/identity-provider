FROM golang:1.14-stretch as base

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/provider

# final image
FROM scratch

WORKDIR /go

# we need this since the binary has a dependency to C libraries (librdkafka)
# COPY --from=base /lib /lib
# COPY --from=base /lib64 /lib64

COPY --from=base /build/provider/templates ./provider/templates
COPY --from=base /build/bin/provider ./bin/provider

ENTRYPOINT ["./bin/provider"]