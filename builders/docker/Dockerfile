FROM golang:1.19-alpine AS src

WORKDIR $GOPATH/src/app/
COPY . .

# Build binary with
# disabling cross compilation, use linux
# removing debug information
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/items_producer ./cmd/itemsproducer/main.go;

RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/items_processor ./cmd/itemsprocessor/main.go;

FROM scratch

#RUN mkdir -p database/migrations # TODO later
#COPY --from=src /go/src/app/migrations/ migrations/
COPY --from=src /go/bin/items_producer .
COPY --from=src /go/bin/items_processor .
COPY ./builders/docker/.env .env

EXPOSE 8080
