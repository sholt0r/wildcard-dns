# wildcard-dns

A simple wildcard DNS server written in Go.
Resolves all `*.<zone>` requests to a user-defined IP (e.g. Traefik) and forwards all other requests to a configurable upstream DNS server.

I built it for me, so if it doesn't work for you, fork it and change it, idc.

## Usage

### Build Locally

```bash
go mod tidy
go build -o wildcard-dns
./wildcard-dns
```

### Docker

```bash
docker build -t wildcard-dns .

docker run --rm --expose 53/udp
-e DNS_PORT=":53"
-e TRAEFIK_IP="172.16.0.201"
-e DOMAIN_ZONE="localhost"
-e UPSTREAM_DNS="127.0.0.11:53"
wildcard-dns
```

### Docker Compose Example

```yaml
service:
  wildcard-dns:
    container_name: wildcard-dns
    image: ghcr.io/sholt0r/wildcard-dns:latest
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
