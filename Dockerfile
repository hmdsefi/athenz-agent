
# Build Athenz-Agent binaries
FROM golang:1.13.6-alpine AS builder

#RUN apk add --update --no-cache ca-certificates make git curl mercurial bzr unzip

WORKDIR /athenz-agent

# Making sure that dependency is not touched
ENV GOFLAGS="-mod=readonly"

# Copy go mod dependencies and build cache
COPY go.* ./
RUN go mod download

# Download dependencies
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o agent cmd/server/main.go

# athenz-agent server
FROM alpine AS athenz-agent-server

COPY --from=builder /athenz-agent/agent /usr/local/bin

EXPOSE 9091
ENTRYPOINT ["agent", "start"]