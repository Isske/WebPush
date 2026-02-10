package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"webpush/database"
	"webpush/models"
	"webpush/utils"
)

var (
	LatestSubscription *models.Subscription
)

// HandleSubscribe processes new push subscription requests
func HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var sub models.Subscription
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, "Invalid subscription", http.StatusBadRequest)
		return
	}

	// Collect IP address
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	} else {
		// X-Forwarded-For may contain multiple IPs
		ip = strings.Split(ip, ",")[0]
	}
	sub.IP = ip

	// Collect User-Agent
	ua := r.Header.Get("User-Agent")
	sub.OS, sub.OSVersion, sub.Browser, sub.BrowserVersion = utils.ParseUserAgent(ua)

	// Use platformVersion from userAgentData if available (more accurate for Windows)
	if sub.Platform != "" && sub.PlatformVersion != "" {
		sub.OS = sub.Platform
		sub.OSVersion = utils.ParsePlatformVersion(sub.Platform, sub.PlatformVersion)
	}

	// Collect nation (GeoIP lookup)
	nation := utils.LookupNation(ip)
	sub.Nation = nation

	LatestSubscription = &sub

	// Save subscription to database
	err = database.SaveSubscription(&sub)
	if err != nil {
		log.Printf("Error saving subscription: %v", err)
		http.Error(w, "Failed to save subscription", http.StatusInternalServerError)
		return
	}

	log.Printf("Subscription received: %s | IP: %s | Nation: %s | OS: %s %s | Browser: %s %s\n",
		sub.Endpoint, sub.IP, sub.Nation, sub.OS, sub.OSVersion, sub.Browser, sub.BrowserVersion)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// LoadSubscriptions reads all subscriptions from database
func LoadSubscriptions() []models.Subscription {
	subs, err := database.GetAllSubscriptions()
	if err != nil {
		log.Printf("Error loading subscriptions: %v", err)
		return []models.Subscription{}
	}
	return subs
}

// SaveSubscriptions saves valid subscriptions to database
func SaveSubscriptions(subs []models.Subscription) error {
	// First, get all current subscriptions
	currentSubs, err := database.GetAllSubscriptions()
	if err != nil {
		log.Printf("Error getting current subscriptions: %v", err)
		return err
	}

	// Create a map of valid endpoints
	validEndpoints := make(map[string]bool)
	for _, sub := range subs {
		validEndpoints[sub.Endpoint] = true
	}

	// Remove subscriptions that are no longer valid
	for _, currentSub := range currentSubs {
		if !validEndpoints[currentSub.Endpoint] {
			database.RemoveSubscription(currentSub.Endpoint)
		}
	}

	return nil
}

// RemoveLatestSubscription removes the latest subscription
func RemoveLatestSubscription() {
	if LatestSubscription != nil {
		database.RemoveSubscription(LatestSubscription.Endpoint)
	}
}
