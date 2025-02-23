package provider

import (
	"context"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/umputun/reproxy/app/discovery"
)

// Static provider, rules are server,from,to
type Static struct {
	Rules []string // each rule is 4 elements comma separated - server,source_url,destination,ping
}

// Events returns channel updating once
func (s *Static) Events(_ context.Context) <-chan discovery.ProviderID {
	res := make(chan discovery.ProviderID, 1)
	res <- discovery.PIStatic
	return res
}

// List all src dst pairs
func (s *Static) List() (res []discovery.URLMapper, err error) {

	parse := func(inp string) (discovery.URLMapper, error) {
		elems := strings.Split(inp, ",")
		if len(elems) != 4 {
			return discovery.URLMapper{}, errors.Errorf("invalid rule %q", inp)
		}
		rx, err := regexp.Compile(strings.TrimSpace(elems[1]))
		if err != nil {
			return discovery.URLMapper{}, errors.Wrapf(err, "can't parse regex %s", elems[1])
		}

		return discovery.URLMapper{
			Server:     strings.TrimSpace(elems[0]),
			SrcMatch:   *rx,
			Dst:        strings.TrimSpace(elems[2]),
			PingURL:    strings.TrimSpace(elems[3]),
			ProviderID: discovery.PIStatic,
		}, nil
	}

	for _, r := range s.Rules {
		if strings.TrimSpace(r) == "" {
			continue
		}
		um, err := parse(r)
		if err != nil {
			return nil, err
		}
		res = append(res, um)
	}
	return res, nil
}
