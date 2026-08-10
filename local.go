package ip2proxycaddy

import (
	"github.com/ip2location/ip2proxy-go/v4"
)

func (h *IP2ProxyCaddy) initLocal() error {
	// Open IP2Proxy BIN database here.
	db, err := ip2proxy.OpenDB(h.BinPath)
	if err != nil {
		return err
	}
	h.db = db

	return nil
}

func (h *IP2ProxyCaddy) lookupLocal(ip string) (ip2proxy.IP2ProxyRecord, error) {
	// Query local BIN here.
	results, err := h.db.GetAll(ip)

	if err != nil {
		return results, err
	}
	return results, nil
}
