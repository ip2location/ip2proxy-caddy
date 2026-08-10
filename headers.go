package ip2proxycaddy

import (
	"github.com/ip2location/ip2location-io-go/ip2locationio"
	"github.com/ip2location/ip2proxy-go/v4"
	"net/http"
	"strconv"
)

func (h *IP2ProxyCaddy) setDBGeoHeaders(r *http.Request, geo ip2proxy.IP2ProxyRecord) {
	prefix := h.HeaderPrefix

	// need to match the logic from the IP2Location.io API
	isProxy := geo.CountryShort != "-" &&
		geo.ProxyType != "AIC" &&
		geo.ProxyType != "SES" &&
		geo.ProxyType != "DCH" &&
		geo.ProxyType != "CDN"

	if isProxy {
		r.Header.Set(prefix+"-Is-Proxy", "True")
	} else {
		r.Header.Set(prefix+"-Is-Proxy", "False")
	}
	r.Header.Set(prefix+"-Proxy-Type", geo.ProxyType)
	r.Header.Set(prefix+"-Country-Code", geo.CountryShort)
	r.Header.Set(prefix+"-Country-Name", geo.CountryLong)
	r.Header.Set(prefix+"-Region", geo.Region)
	r.Header.Set(prefix+"-City", geo.City)
	r.Header.Set(prefix+"-ISP", geo.Isp)
	r.Header.Set(prefix+"-Domain", geo.Domain)
	r.Header.Set(prefix+"-Usage-Type", geo.UsageType)
	r.Header.Set(prefix+"-ASN", geo.Asn)
	r.Header.Set(prefix+"-AS", geo.As)
	r.Header.Set(prefix+"-Last-Seen", geo.LastSeen)
	r.Header.Set(prefix+"-Threat", geo.Threat)
	r.Header.Set(prefix+"-Provider", geo.Provider)
	r.Header.Set(prefix+"-Fraud-Score", geo.FraudScore)
}

func (h *IP2ProxyCaddy) setAPIGeoHeaders(r *http.Request, geo ip2locationio.IPGeolocationResult) {
	prefix := h.HeaderPrefix

	// currently, we are matching the fields returned by the DB case
	if geo.IsProxy {
		r.Header.Set(prefix+"-Is-Proxy", "True")
	} else {
		r.Header.Set(prefix+"-Is-Proxy", "False")
	}
	r.Header.Set(prefix+"-Proxy-Type", geo.Proxy.ProxyType)
	r.Header.Set(prefix+"-Country-Code", geo.CountryCode)
	r.Header.Set(prefix+"-Country-Name", geo.CountryName)
	r.Header.Set(prefix+"-Region", geo.RegionName)
	r.Header.Set(prefix+"-City", geo.CityName)
	r.Header.Set(prefix+"-ISP", geo.Isp)
	r.Header.Set(prefix+"-Domain", geo.Domain)
	r.Header.Set(prefix+"-Usage-Type", geo.UsageType)
	r.Header.Set(prefix+"-ASN", geo.Asn)
	r.Header.Set(prefix+"-AS", geo.AS)
	r.Header.Set(prefix+"-Last-Seen", strconv.Itoa(geo.Proxy.LastSeen))
	r.Header.Set(prefix+"-Threat", geo.Proxy.Threat)
	r.Header.Set(prefix+"-Provider", geo.Proxy.Provider)
	r.Header.Set(prefix+"-Fraud-Score", strconv.Itoa(geo.FraudScore))
}
