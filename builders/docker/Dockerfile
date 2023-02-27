FROM golang:1.19-alpine AS src

WORKDIR $GOPATH/src/app/
COPY . .

# Build binary with
# disabling cross compilation
# removing debug information
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/server ./cmd/sample/main.go;

RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/cover ./cmd/coverage/main.go;

FROM scratch

#RUN mkdir -p database/migrations #if fs exists
COPY --from=src /go/src/app/migrations/ migrations/
COPY --from=src /go/bin/server .
COPY --from=src /go/bin/cover .
COPY ./builders/docker/.env .env

EXPOSE 8080
