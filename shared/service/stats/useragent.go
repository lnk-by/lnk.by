package stats

import (
	"fmt"

	"github.com/mileusna/useragent"
)

func useragentParser(e Event) string {
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
		incr(ua.Desktop), incr(ua.Tablet), incr(ua.Mobile), incr(ua.Bot), incr(!(ua.Desktop || ua.Tablet || ua.Mobile || ua.Bot)),
		// OS
		incr(ua.IsWindows()), incr(ua.IsLinux()), incr(ua.IsMacOS()), incr(ua.IsIOS()), incr(ua.IsAndroid()),
		incr(!(ua.IsWindows() || ua.IsLinux() || ua.IsMacOS() || ua.IsIOS() || ua.IsAndroid())),
		// Browser
		incr(ua.IsChrome()), incr(ua.IsOpera() || ua.IsOperaMini()), incr(ua.IsInternetExplorer()), incr(ua.IsEdge()), incr(ua.IsFirefox()),
		incr(!(ua.IsChrome() || ua.IsOpera() || ua.IsOperaMini() || ua.IsInternetExplorer() || ua.IsEdge() || ua.IsFirefox())),

		incr(ua.Desktop && ua.IsWindows()), incr(ua.Desktop && ua.IsLinux()), incr(ua.Desktop && ua.IsMacOS()),
		incr(ua.Tablet && ua.IsWindows()), incr(ua.Tablet && ua.IsLinux()), incr(ua.Tablet && ua.IsIOS()),
		incr(ua.Mobile && ua.IsAndroid()), incr(ua.Mobile && ua.IsIOS()),
	)
}

func incr(f bool) int {
	if f {
		return 1
	}
	return 0
}
