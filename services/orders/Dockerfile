FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/orders ./services/orders

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /out/orders /app/orders

EXPOSE 8080 8081

USER nonroot:nonroot

ENTRYPOINT ["/app/orders"]