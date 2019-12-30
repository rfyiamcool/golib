package dnsquery

import (
	"errors"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/miekg/dns"
)

type DnsResolver struct {
	Servers    []string
	RetryTimes int

	r *rand.Rand
}

func New(servers []string) *DnsResolver {
	for i := range servers {
		servers[i] = net.JoinHostPort(servers[i], "53")
	}

	return &DnsResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func NewFromResolvConf(path string) (*DnsResolver, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &DnsResolver{}, errors.New("no such file or directory: " + path)
	}

	config, err := dns.ClientConfigFromFile(path)
	servers := []string{}
	for _, ipAddress := range config.Servers {
		servers = append(servers, net.JoinHostPort(ipAddress, "53"))
	}

	return &DnsResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}, err
}

func (r *DnsResolver) LookupHost(host string) ([]net.IP, error) {
	var (
		res []net.IP
		err error
	)

	for r.RetryTimes > 0 {
		res, err = r.lookupHost(host)
		if err == nil {
			return res, err
		}

		r.RetryTimes--
	}

	return res, err
}

func (r *DnsResolver) lookupHost(host string) ([]net.IP, error) {
	m := new(dns.Msg)
	m.Id = dns.Id()
	m.RecursionDesired = true
	m.SetQuestion(dns.Fqdn(host), dns.TypeA)

	// query
	in, err := dns.Exchange(m, r.Servers[r.r.Intn(len(r.Servers))])

	result := []net.IP{}
	if err != nil {
		return result, err
	}

	if in != nil && in.Rcode != dns.RcodeSuccess {
		return result, errors.New(dns.RcodeToString[in.Rcode])
	}

	for _, record := range in.Answer {
		if t, ok := record.(*dns.A); ok {
			result = append(result, t.A)
		}
	}

	return result, err
}
