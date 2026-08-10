package ip2proxycaddy

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/ip2location/ip2location-io-go/ip2locationio"
	"github.com/ip2location/ip2proxy-go/v4"
	"net/http"
)

func init() {
	caddy.RegisterModule(IP2ProxyCaddy{})
}

// IP2Proxy is a Caddy HTTP handler that detects the visitor IP address,
// looks up proxy data using either a local IP2Proxy BIN database
// or the IP2Location.io API, then injects the result as request headers
// for upstream handlers such as reverse_proxy or php_fastcgi.
type IP2ProxyCaddy struct {
	// Mode specifies the lookup method.
	// Valid values are "local" for BIN database lookup or "remote" for API lookup.
	Mode string `json:"mode,omitempty"`

	// BinPath is the filesystem path to the IP2Proxy BIN database.
	// Required when mode is "local".
	BinPath string `json:"bin_path,omitempty"`

	// APIKey is the IP2Location.io API key.
	// Required when mode is "remote".
	APIKey string `json:"api_key,omitempty"`

	// HeaderPrefix is the prefix used for injected headers.
	// Defaults to "X-IP2Proxy".
	HeaderPrefix string `json:"header_prefix,omitempty"`

	APIURL string
	db     *ip2proxy.DB
	api    *ip2locationio.IPGeolocation
}

func (IP2ProxyCaddy) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.ip2proxy",
		New: func() caddy.Module { return new(IP2ProxyCaddy) },
	}
}

func (h *IP2ProxyCaddy) Provision(ctx caddy.Context) error {
	if h.Mode == "" {
		h.Mode = "local"
	}

	if h.HeaderPrefix == "" {
		h.HeaderPrefix = "X-IP2Proxy"
	}

	switch h.Mode {
	case "local":
		if h.BinPath == "" {
			return fmt.Errorf("bin_path is required when mode is local")
		}
		return h.initLocal()

	case "remote":
		h.APIURL = "https://api.ip2location.io"
		if h.APIKey == "" {
			return fmt.Errorf("api_key is required when mode is remote")
		}
		return h.initRemote()

	default:
		return fmt.Errorf("mode must be local or remote")
	}
}

func (h *IP2ProxyCaddy) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	ip := getClientIP(r)

	var geoDB ip2proxy.IP2ProxyRecord
	var geoAPI ip2locationio.IPGeolocationResult
	var err error

	switch h.Mode {
	case "local":
		geoDB, err = h.lookupLocal(ip)
		if err == nil {
			h.setDBGeoHeaders(r, geoDB)
		}
	case "remote":
		geoAPI, err = h.lookupRemote(ip)
		if err == nil {
			h.setAPIGeoHeaders(r, geoAPI)
		}
	}

	return next.ServeHTTP(w, r)
}

var (
	_ caddy.Provisioner           = (*IP2ProxyCaddy)(nil)
	_ caddyhttp.MiddlewareHandler = (*IP2ProxyCaddy)(nil)
)
