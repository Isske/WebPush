package handlers

import (
	"encoding/json"
	"net/http"
	"webpush/database"
	"webpush/models"
)

// GetDashboardStatsHandler returns statistics for the dashboard
func GetDashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	subs := LoadSubscriptions()

	// Get countries from database with proper aggregation
	countries, err := database.GetCountByNation()
	if err != nil {
		countries = make(map[string]int)
	}

	stats := models.DashboardStats{
		TotalClients:  len(subs),
		OnlineClients: len(subs), // For now, all subscriptions are considered online
		TotalPushes:   database.GetPushCount(),
		Countries:     countries,
		Subscriptions: subs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ServeDashboard serves the dashboard HTML page
func ServeDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}
