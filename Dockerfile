# minimal go build
FROM golang:1.14 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
# cache
RUN go mod download

COPY cmd/ cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o kafka_sink cmd/kafka_sink/main.go

# https://github.com/GoogleContainerTools/distroless 
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kafka_sink .
USER nonroot:nonroot

ENTRYPOINT ["/kafka_sink"]
