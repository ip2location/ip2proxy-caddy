package ip2proxycaddy

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("ip2proxy", parseCaddyfile)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m IP2ProxyCaddy
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return &m, err
}

func (h *IP2ProxyCaddy) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "mode":
				if !d.NextArg() {
					return d.ArgErr()
				}
				h.Mode = d.Val()

			case "bin_path":
				if !d.NextArg() {
					return d.ArgErr()
				}
				h.BinPath = d.Val()

			case "api_key":
				if !d.NextArg() {
					return d.ArgErr()
				}
				h.APIKey = d.Val()

			case "header_prefix":
				if !d.NextArg() {
					return d.ArgErr()
				}
				h.HeaderPrefix = d.Val()

			default:
				return d.Errf("unknown ip2proxy option: %s", d.Val())
			}
		}
	}

	return nil
}

var _ caddyfile.Unmarshaler = (*IP2ProxyCaddy)(nil)
