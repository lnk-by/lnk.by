package maxmind

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"os"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var (
	mmdb    *geoip2.Reader
	once    sync.Once
	initErr error
)

const mmdbFile = "GeoLite2-Country.mmdb"

func Init() error {
	once.Do(func() {
		if _, err := os.Stat(mmdbFile); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				initErr = errors.New("maxmind database file not found")
			} else {
				initErr = err // unknown file access error
			}
			return
		}

		// Try to open the DB
		mmdb, initErr = geoip2.Open(mmdbFile)
	})
	return initErr
}

func IPToCountry(ctx context.Context, ip string) string {
	if mmdb == nil {
		return ""
	}
	parsedIP := net.ParseIP(ip)
	record, err := mmdb.Country(parsedIP)
	if err != nil || record.Country.IsoCode == "" {
		return ""
	}
	return record.Country.IsoCode
}
