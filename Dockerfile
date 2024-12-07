# Build
FROM mirror.gcr.io/library/golang:1.23 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./

ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /vzor ./cmd/vzor

# Tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy
FROM mirror.gcr.io/library/debian:11-slim AS build-release-stage

WORKDIR /

COPY --from=build-stage /vzor /vzor
RUN apt-get update
RUN apt-get install -y ca-certificates

EXPOSE 8080

ENTRYPOINT [ "/vzor"]