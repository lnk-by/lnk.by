package stats

import (
	"fmt"

	"github.com/mileusna/useragent"
)

func updateUserAgentBasedStatistics(e Event) string {
	ua := useragent.Parse(e.UserAgent)
	return fmt.Sprintf(`
			UPDATE useragent_count SET 
				device_desktop = device_desktop + %d,
				device_tablet = device_tablet + %d,
				device_mobile = device_mobile + %d,
				device_bot = device_bot + %d,
				device_other = device_other + %d,

				os_windows = os_windows + %d,
				os_linux = os_linux + %d,
				os_macos = os_macos + %d,
				os_ios = os_ios + %d,
				os_android = os_android + %d,
				os_other = os_other + %d,

				browser_chrome = browser_chrome + %d,
				browser_opera = browser_opera + %d,
				browser_internet_exporer = browser_internet_exporer + %d,
				browser_edge = browser_edge + %d,
				browser_firefox = browser_firefox + %d,
				browser_other = browser_other + %d,

				desktop_windows = desktop_windows + %d,
				desktop_linux = desktop_linux + %d,
				desktop_macos = desktop_macos + %d,
				tablet_windows = tablet_windows + %d,
				tablet_linux = tablet_linux + %d,
				tablet_ios = tablet_ios + %d,

				mobile_android = mobile_android + %d,
				mobile_ios = mobile_ios + %d
			WHERE key = $1`,
		// device
		b2i(ua.Desktop), b2i(ua.Tablet), b2i(ua.Mobile), b2i(ua.Bot), b2i(!(ua.Desktop || ua.Tablet || ua.Mobile || ua.Bot)),
		// OS
		b2i(ua.IsWindows()), b2i(ua.IsLinux()), b2i(ua.IsMacOS()), b2i(ua.IsIOS()), b2i(ua.IsAndroid()),
		b2i(!(ua.IsWindows() || ua.IsLinux() || ua.IsMacOS() || ua.IsIOS() || ua.IsAndroid())),
		// Browser
		b2i(ua.IsChrome()), b2i(ua.IsOpera() || ua.IsOperaMini()), b2i(ua.IsInternetExplorer()), b2i(ua.IsEdge()), b2i(ua.IsFirefox()),
		b2i(!(ua.IsChrome() || ua.IsOpera() || ua.IsOperaMini() || ua.IsInternetExplorer() || ua.IsEdge() || ua.IsFirefox())),

		b2i(ua.Desktop && ua.IsWindows()), b2i(ua.Desktop && ua.IsLinux()), b2i(ua.Desktop && ua.IsMacOS()),
		b2i(ua.Tablet && ua.IsWindows()), b2i(ua.Tablet && ua.IsLinux()), b2i(ua.Tablet && ua.IsIOS()),
		b2i(ua.Mobile && ua.IsAndroid()), b2i(ua.Mobile && ua.IsIOS()),
	)
}

// b2i translates boolean to int: true->1, false->0
func b2i(f bool) int {
	if f {
		return 1
	}
	return 0
}
