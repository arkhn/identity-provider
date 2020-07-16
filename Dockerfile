FROM golang:1.14-stretch as base

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN chmod +x ./wait-for-postgres.sh

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/provider

# final image
FROM postgres

WORKDIR /go

COPY --from=base /build/wait-for-postgres.sh ./wait-for-postgres.sh

COPY --from=base /build/provider/templates ./provider/templates
COPY --from=base /build/bin/provider ./bin/provider

ENTRYPOINT ["./wait-for-postgres.sh", "./bin/provider"]