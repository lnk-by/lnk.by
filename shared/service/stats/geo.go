package stats

import (
	"context"
	"fmt"
)

func geoFactory(ipToLocFn func(context.Context, string) string) func(context.Context, Event) string {
	return func(ctx context.Context, e Event) string {
		country := ipToLocFn(ctx, e.IP)
		if country == "" {
			country = "UNKNOWN"
		}
		country = "C_" + country
		return fmt.Sprintf("UPDATE country_count SET %[1]s = %[1]s + 1 WHERE key = $1", country)
	}
}
