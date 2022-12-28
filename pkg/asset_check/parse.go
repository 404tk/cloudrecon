package asset_check

import (
	"fmt"
	"strings"

	"github.com/404tk/cloudrecon/common"
	"golang.org/x/net/publicsuffix"
)

type DomainFormat struct {
	Subdomain, Domain, TLD string
	ICANN                  bool
}

// Code source: https://github.com/jpillora/go-tld/blob/master/parse.go#L24
func ParseDomain(s string) (*DomainFormat, error) {
	etld1, err := publicsuffix.EffectiveTLDPlusOne(s)
	suffix, icann := publicsuffix.PublicSuffix(strings.ToLower(s))
	// HACK: attempt to support valid domains which are not registered with ICAN
	if err != nil && !icann && suffix == s {
		etld1 = s
		err = nil
	}
	if err != nil {
		return nil, err
	}
	//convert to domain name, and tld
	i := strings.Index(etld1, ".")
	if i < 0 {
		return nil, fmt.Errorf("tld: failed parsing %q", s)
	}
	domName := etld1[0:i]
	tld := etld1[i+1:]
	//and subdomain
	sub := ""
	if rest := strings.TrimSuffix(s, "."+etld1); rest != s {
		sub = rest
	}
	return &DomainFormat{
		Subdomain: sub,
		Domain:    domName,
		TLD:       tld,
		ICANN:     icann,
	}, nil
}

type Payload struct {
	Urls   map[string]string
	Method string
	Exp    []int
}

func ParsePayloads(words []string) ([]Payload, int) {
	data := []Payload{}
	num := 0
	for _, word := range words {
		for _, poc := range defaultPocs {
			for _, rule := range poc.Rules {
				urls := ParseUrls(word, rule)
				num += len(urls)
				data = append(data, Payload{
					Urls:   urls,
					Method: rule.Method,
					Exp:    rule.Expression})
			}
		}
	}
	return data, num
}

func ParseUrls(word string, r Rule) map[string]string {
	urls := make(map[string]string)
	for _, sep := range common.Separators {
		for _, env := range common.Environments {
			for _, format := range r.Format {
				if !((sep == "-" || sep == "_") && env == "") {
					url := strings.Replace(format, "{word}", word, -1)
					url = strings.Replace(url, "{sep}", sep, -1)
					url = strings.Replace(url, "{env}", env, -1)
					if r.Traversal {
						for _, region := range r.Regions {
							_url := strings.Replace(url, "{region}", region, -1)
							urls[_url] = region
						}
					} else {
						url = strings.Replace(url, "{region}", r.Regions[0], -1)
						urls[url] = r.Regions[0]
					}
				}
			}
		}
	}

	return urls
}
