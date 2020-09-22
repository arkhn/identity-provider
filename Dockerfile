FROM golang:1.14-stretch as base

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN chmod +x ./wait-for-postgres.sh

# Build identity provider
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/provider
# Build superuser script
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/superuser scripts/superuser.go

# final image
FROM postgres:13-alpine

WORKDIR /go
ENV PATH=$PATH:/go/bin

COPY --from=base /build/wait-for-postgres.sh ./wait-for-postgres.sh

COPY --from=base /build/provider/templates ./provider/templates
COPY --from=base /build/bin/provider ./bin/provider
COPY --from=base /build/bin/superuser ./bin/seed-superuser

ENTRYPOINT ["./wait-for-postgres.sh", "./bin/provider"]