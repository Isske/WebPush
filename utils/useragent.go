package utils

import (
	"strconv"
	"strings"
)

// ParseUserAgent extracts OS and browser information from a user agent string
// Returns: (os, osVersion, browser, browserVersion)
func ParseUserAgent(ua string) (string, string, string, string) {
	os := "Unknown"
	osVersion := ""
	browser := "Unknown"
	browserVersion := ""

	// Extract OS and version
	if strings.Contains(ua, "Windows NT") {
		os = "Windows"
		if idx := strings.Index(ua, "Windows NT "); idx != -1 {
			start := idx + 11
			end := strings.IndexAny(ua[start:], ";) ")
			if end != -1 {
				ntVersion := ua[start : start+end]
				switch ntVersion {
				case "10.0":
					osVersion = "10/11"
				case "6.3":
					osVersion = "8.1"
				case "6.2":
					osVersion = "8"
				case "6.1":
					osVersion = "7"
				default:
					osVersion = ntVersion
				}
			}
		}
	} else if strings.Contains(ua, "Mac OS X") {
		os = "macOS"
		if idx := strings.Index(ua, "Mac OS X "); idx != -1 {
			start := idx + 9
			end := strings.IndexAny(ua[start:], ";)")
			if end != -1 {
				osVersion = strings.ReplaceAll(ua[start:start+end], "_", ".")
			}
		}
	} else if strings.Contains(ua, "Android") {
		os = "Android"
		if idx := strings.Index(ua, "Android "); idx != -1 {
			start := idx + 8
			end := strings.IndexAny(ua[start:], ";)")
			if end != -1 {
				osVersion = ua[start : start+end]
			}
		}
	} else if strings.Contains(ua, "iPhone OS") || strings.Contains(ua, "CPU OS") {
		os = "iOS"
		var idx int
		if strings.Contains(ua, "iPhone OS") {
			idx = strings.Index(ua, "iPhone OS ")
			if idx != -1 {
				start := idx + 10
				end := strings.IndexAny(ua[start:], " )")
				if end != -1 {
					osVersion = strings.ReplaceAll(ua[start:start+end], "_", ".")
				}
			}
		} else if strings.Contains(ua, "CPU OS") {
			idx = strings.Index(ua, "CPU OS ")
			if idx != -1 {
				start := idx + 7
				end := strings.IndexAny(ua[start:], " )")
				if end != -1 {
					osVersion = strings.ReplaceAll(ua[start:start+end], "_", ".")
				}
			}
		}
	} else if strings.Contains(ua, "iPad") {
		os = "iPadOS"
		if idx := strings.Index(ua, "CPU OS "); idx != -1 {
			start := idx + 7
			end := strings.IndexAny(ua[start:], " )")
			if end != -1 {
				osVersion = strings.ReplaceAll(ua[start:start+end], "_", ".")
			}
		}
	} else if strings.Contains(ua, "Linux") {
		os = "Linux"
	}

	// Extract browser and version
	if strings.Contains(ua, "Edg/") {
		browser = "Edge"
		if idx := strings.Index(ua, "Edg/"); idx != -1 {
			start := idx + 4
			end := strings.IndexAny(ua[start:], " )")
			if end != -1 {
				browserVersion = ua[start : start+end]
			}
		}
	} else if strings.Contains(ua, "Chrome/") {
		browser = "Chrome"
		if idx := strings.Index(ua, "Chrome/"); idx != -1 {
			start := idx + 7
			end := strings.IndexAny(ua[start:], " )")
			if end != -1 {
				browserVersion = ua[start : start+end]
			}
		}
	} else if strings.Contains(ua, "Firefox/") {
		browser = "Firefox"
		if idx := strings.Index(ua, "Firefox/"); idx != -1 {
			start := idx + 8
			end := strings.IndexAny(ua[start:], " )")
			if end != -1 {
				browserVersion = ua[start : start+end]
			}
		}
	} else if strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/") {
		browser = "Safari"
		if idx := strings.Index(ua, "Version/"); idx != -1 {
			start := idx + 8
			end := strings.IndexAny(ua[start:], " )")
			if end != -1 {
				browserVersion = ua[start : start+end]
			}
		}
	} else if strings.Contains(ua, "OPR/") {
		browser = "Opera"
		if idx := strings.Index(ua, "OPR/"); idx != -1 {
			start := idx + 4
			end := strings.IndexAny(ua[start:], " )")
			if end != -1 {
				browserVersion = ua[start : start+end]
			}
		}
	}

	return os, osVersion, browser, browserVersion
}

// ParsePlatformVersion converts platform version from userAgentData to human-readable OS version
func ParsePlatformVersion(platform, platformVersion string) string {
	if platform == "Windows" {
		// Parse major version from platformVersion (e.g., "13.0.0" -> 13)
		parts := strings.Split(platformVersion, ".")
		if len(parts) > 0 {
			majorVersion := parts[0]
			// Convert string to int to compare
			if major, err := strconv.Atoi(majorVersion); err == nil {
				if major >= 13 {
					return "11"
				} else if major >= 10 {
					return "10"
				} else if major >= 6 {
					// Older Windows versions
					switch majorVersion {
					case "6":
						if len(parts) > 1 {
							minor := parts[1]
							switch minor {
							case "3":
								return "8.1"
							case "2":
								return "8"
							case "1":
								return "7"
							}
						}
					}
				}
			}
		}
		return platformVersion
	}
	// For other platforms, return as-is
	return platformVersion
}
