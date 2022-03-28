FROM golang:1.18 as builder
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go vet ./...
ARG version=0.1.0-alpha buildDate=2022-03-27
RUN CGO_ENABLED=0 go build -a \
    -ldflags "-X main.version=$version -X main.buildDate=$buildDate" \
    -o bin/blueprint \
    -installsuffix static \
    ./cmd/

FROM gcr.io/distroless/static:nonroot
FROM busybox
COPY --from=builder /workspace/bin/blueprint /
ENTRYPOINT [ "/blueprint" ]
