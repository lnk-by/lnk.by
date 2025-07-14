package stats

import (
	"context"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var mmdb *geoip2.Reader

func init() {
	var err error
	mmdb, err = geoip2.Open("GeoLite2-Country.mmdb") // in deployment package
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
}

func ipToCountry(ctx context.Context, ip string) string {
	parsedIP := net.ParseIP(ip)
	record, err := mmdb.Country(parsedIP)
	if err != nil || record.Country.IsoCode == "" {
		return ""
	}
	return record.Country.IsoCode
}
