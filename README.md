# wildcarddns

A simple wildcard DNS server written in Go.
Resolves all `*.localhost` domains to a user-defined IP (e.g. Traefik) and forwards all other requests to a configurable upstream DNS server.

## Features

* Wildcard resolution for `*.localhost`
* Upstream forwarding for other domains
* Lightweight, fast, and container-ready
* Configuration via environment variables

## Usage

### Build Locally

```bash
go mod tidy
go build -o dnsserver
./dnsserver
```

### Docker

```bash
docker build -t wildcarddns .
docker run --rm --expose 53/udp
-e DNS\_PORT=":53"
-e TRAEFIK\_IP="172.19.0.250"
-e DOMAIN\_ZONE="localhost"
-e UPSTREAM\_DNS="1.1.1.1:53"
\--cap-add=NET\_ADMIN
wildcarddns
```

### Docker Compose Example

```yaml
service:
  wildcard-dns:
    container_name: wildcard-dns
    image: sholt0r/wildcard-dns:latest
    restart: unless-stopped
    networks:
        proxy:
            ipv4_address: 172.16.0.200
    expose:
        - 53/udp
    environment:
        DNS_PORT:       :53
        DNS_PROXY:      172.16.0.201
        DNS_ZONE:       localhost
        DNS_UPSTREAM:   127.0.0.11

networks:
    proxy:
        name: proxy
        external: true

```

