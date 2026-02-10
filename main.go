package main

import (
	"log"
	"net/http"
	"webpush/database"
	"webpush/handlers"
)

func main() {
	// Initialize database
	if err := database.InitDB("data/webpush.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized")

	// Initialize VAPID keys
	if err := handlers.InitVAPIDKeys(); err != nil {
		log.Fatalf("Failed to initialize VAPID keys: %v", err)
	}

	// Setup HTTP routes
	// Dashboard routes
	http.HandleFunc("/", handlers.ServeDashboard)
	http.HandleFunc("/api/stats", handlers.GetDashboardStatsHandler)

	// Push notification routes
	http.HandleFunc("/vapid-public-key", handlers.GetVAPIDPublicKeyHandler)
	http.HandleFunc("/subscribe", handlers.HandleSubscribe)
	http.HandleFunc("/send-notification", handlers.SendNotificationHandler)
	http.HandleFunc("/send-broadcast", handlers.SendBroadcastHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve service worker from root
	http.Handle("/sw.js", http.FileServer(http.Dir("static")))

	port := ":10040"
	log.Printf("WebPush Server starting on http://localhost%s\n", port)
	log.Printf("Dashboard available at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
