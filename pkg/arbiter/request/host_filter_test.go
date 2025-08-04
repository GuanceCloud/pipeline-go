package request

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHostFilte(t *testing.T) {
	cases := []struct {
		name         string
		blockdCIDRs  []string
		allowedCIDRs []string
		blockdHosts  []string
		allowedHosts []string

		addr    string
		blocked bool
	}{
		{
			name: "test",
			addr: "baidu.com:80",
		},
		{
			name:        "test1",
			addr:        "baidu.com:80",
			blockdHosts: []string{"abc.com", "baidu.com", "qq.com"},
			blocked:     true,
		},
		{
			name:         "test2",
			addr:         "baidu.com:80",
			allowedHosts: []string{"abc.com", "qq.com"},
			blocked:      true,
		},
		{
			name:         "test3",
			addr:         "baidu.com:80",
			allowedCIDRs: []string{"1.0.0.0/8"},
			blocked:      true,
		},
		{
			name:         "test4",
			addr:         "dns.alidns.com:80",
			allowedCIDRs: []string{"223.0.0.0/8", "2400::/16"},
		},
		{
			name:        "test5",
			addr:        "local.ubwbu.com:1234",
			blockdCIDRs: PrivateCIDRs(),
			blocked:     true,
		},
		{
			name:        "test6",
			addr:        "127.0.0.1:80",
			blockdCIDRs: PrivateCIDRs(),
			blocked:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			host, _, err := net.SplitHostPort(c.addr)
			if err != nil {
				t.Fatal(err)
			}

			f := NewHostFilter(c.blockdCIDRs, c.allowedCIDRs, c.blockdHosts, c.allowedHosts, 100, 1*time.Second)

			blocked, _ := f.IsBlocked(host)
			assert.Equal(t, c.blocked, blocked)
		})
	}
}
