package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var (
	dnsPort			= os.Getenv("DNS_PORT")
	proxyIP			= os.Getenv("DNS_PROXY")
	domainZone	= ensureDot(os.Getenv("DNS_ZONE"))
	upstreamDNS	= os.Getenv("DNS_UPSTREAM")
)

func ensureDot(domain string) string {
	if !strings.HasSuffix(domain, ".") {
		return domain + "."
	}
	return domain
}

func handleDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		log.Printf("Query: %s %d", q.Name, q.Qtype)

		if strings.HasSuffix(q.Name, domainZone) && q.Qtype == dns.TypeA {
			switch q.Qtype {
			case dns.TypeA:
				rr := &dns.A{
					Hdr: dns.RR_Header{
						Name: 	q.Name,
						Rrtype:	dns.TypeA,
						Class: 	dns.ClassINET,
						Ttl:		60,
					},
					A: net.ParseIP(proxyIP),
				}
				msg.Answer = append(msg.Answer, rr)
			case dns.TypeAAAA:
				log.Printf("Ignoring AAAA query for %s", q.Name)
			default:
				log.Printf("Unsupported query type %d for %s", q.Qtype, q.Name)
			}
			w.WriteMsg(&msg)
			return
		}
		resp, err := forwardQuery(r)
		if err != nil {
			log.Printf("Forwarding error: %v", err)
			dns.HandleFailed(w, r)
			return
		}
		w.WriteMsg(resp)
	}
}

func forwardQuery(msg *dns.Msg) (*dns.Msg, error) {
	client := &dns.Client{Timeout: 2 * time.Second}
	resp, _, err := client.Exchange(msg, upstreamDNS)
	return resp, err
}

func main() {
	dns.HandleFunc(".", handleDNS)

	server := &dns.Server{Addr: dnsPort, Net: "udp"}
	log.Printf("Starting DNS server on %s for zone *.%s -> %s (forwarding to %s)",
			dnsPort, domainZone, proxyIP, upstreamDNS)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("DNS server error: %v", err)
	}
}
