FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o wildcard-dns

FROM alpine
WORKDIR /app
COPY --from=builder /app/wildcard-dns .
ENTRYPOINT [ "./wildcard-dns" ]
