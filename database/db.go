package database

import (
	"database/sql"
	"log"
	"time"
	"webpush/models"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes the SQLite database
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Create tables
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			endpoint TEXT NOT NULL UNIQUE,
			p256dh TEXT NOT NULL,
			auth TEXT NOT NULL,
			ip TEXT,
			nation TEXT,
			os TEXT,
			os_version TEXT,
			browser TEXT,
			browser_version TEXT,
			platform TEXT,
			platform_version TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_active DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS push_stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			total_pushes INTEGER DEFAULT 0
		);

		-- Initialize push stats if not exists
		INSERT OR IGNORE INTO push_stats (id, total_pushes) VALUES (1, 0);

		CREATE INDEX IF NOT EXISTS idx_endpoint ON subscriptions(endpoint);
		CREATE INDEX IF NOT EXISTS idx_nation ON subscriptions(nation);
	`)

	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// SaveSubscription saves or updates a subscription in the database
func SaveSubscription(sub *models.Subscription) error {
	_, err := DB.Exec(`
		INSERT INTO subscriptions (endpoint, p256dh, auth, ip, nation, os, os_version, browser, browser_version, platform, platform_version)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(endpoint) DO UPDATE SET
			ip = excluded.ip,
			nation = excluded.nation,
			os = excluded.os,
			os_version = excluded.os_version,
			browser = excluded.browser,
			browser_version = excluded.browser_version,
			platform = excluded.platform,
			platform_version = excluded.platform_version,
			last_active = CURRENT_TIMESTAMP
	`, sub.Endpoint, sub.Keys.P256dh, sub.Keys.Auth, sub.IP, sub.Nation, sub.OS, sub.OSVersion, sub.Browser, sub.BrowserVersion, sub.Platform, sub.PlatformVersion)

	return err
}

// GetAllSubscriptions retrieves all subscriptions from the database
func GetAllSubscriptions() ([]models.Subscription, error) {
	rows, err := DB.Query(`
		SELECT endpoint, p256dh, auth, ip, nation, os, os_version, browser, browser_version, platform, platform_version
		FROM subscriptions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.Endpoint,
			&sub.Keys.P256dh,
			&sub.Keys.Auth,
			&sub.IP,
			&sub.Nation,
			&sub.OS,
			&sub.OSVersion,
			&sub.Browser,
			&sub.BrowserVersion,
			&sub.Platform,
			&sub.PlatformVersion,
		)
		if err != nil {
			log.Printf("Error scanning subscription: %v", err)
			continue
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

// RemoveSubscription removes a subscription by endpoint
func RemoveSubscription(endpoint string) error {
	_, err := DB.Exec("DELETE FROM subscriptions WHERE endpoint = ?", endpoint)
	return err
}

// IncrementPushCount increments the total push count
func IncrementPushCount(count int) error {
	_, err := DB.Exec("UPDATE push_stats SET total_pushes = total_pushes + ? WHERE id = 1", count)
	return err
}

// GetPushCount returns the total push count
func GetPushCount() int {
	var count int
	err := DB.QueryRow("SELECT total_pushes FROM push_stats WHERE id = 1").Scan(&count)
	if err != nil {
		log.Printf("Error getting push count: %v", err)
		return 0
	}
	return count
}

// GetCountByNation returns a map of nation to subscriber count
func GetCountByNation() (map[string]int, error) {
	rows, err := DB.Query(`
		SELECT nation, COUNT(*) as count
		FROM subscriptions
		WHERE nation IS NOT NULL AND nation != ''
		GROUP BY nation
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var nation string
		var count int
		if err := rows.Scan(&nation, &count); err != nil {
			continue
		}
		counts[nation] = count
	}

	return counts, nil
}

// GetCountByBrowser returns a map of browser to count
func GetCountByBrowser() (map[string]int, error) {
	rows, err := DB.Query(`
		SELECT browser, COUNT(*) as count
		FROM subscriptions
		WHERE browser IS NOT NULL AND browser != ''
		GROUP BY browser
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var browser string
		var count int
		if err := rows.Scan(&browser, &count); err != nil {
			continue
		}
		counts[browser] = count
	}

	return counts, nil
}

// GetCountByOS returns a map of OS to count
func GetCountByOS() (map[string]int, error) {
	rows, err := DB.Query(`
		SELECT os, COUNT(*) as count
		FROM subscriptions
		WHERE os IS NOT NULL AND os != ''
		GROUP BY os
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var os string
		var count int
		if err := rows.Scan(&os, &count); err != nil {
			continue
		}
		counts[os] = count
	}

	return counts, nil
}

// CleanupOldSubscriptions removes subscriptions inactive for more than the specified duration
func CleanupOldSubscriptions(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	_, err := DB.Exec("DELETE FROM subscriptions WHERE last_active < ?", cutoff)
	return err
}
