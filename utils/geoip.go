package utils

import (
	"io"
	"net/http"
	"strings"
)

// LookupNation returns a country code for an IP using ip-api.com
// Returns empty string on any error (best effort, non-blocking)
func LookupNation(ip string) string {
	if ip == "" {
		return ""
	}

	// Only use the IP part if port is present
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	url := "http://ip-api.com/line/" + ip + "?fields=countryCode"
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
