# WebPush - Unified Push Notification Server

A complete web push notification server with integrated dashboard, built in Go.

## Project Structure

```
webpush/
├── main.go                   # Main entry point
├── go.mod                    # Go module definition
├── handlers/                 # HTTP request handlers
│   ├── vapid.go             # VAPID key management
│   ├── subscription.go      # Subscription handling
│   ├── notification.go      # Push notification sending
│   └── dashboard.go         # Dashboard API
├── models/                   # Data models
│   └── types.go             # Shared types and structures
├── utils/                    # Utility functions
│   ├── useragent.go         # User agent parsing
│   └── geoip.go             # GeoIP lookups
├── static/                   # Web assets
│   ├── index.html           # Dashboard frontend
│   └── sw.js                # Service worker
└── data/                     # Runtime data (generated)
    ├── vapid_public.txt     # VAPID public key
    ├── vapid_private.txt    # VAPID private key
    ├── subscriptions.json   # Subscriber data
    └── push_count.json      # Push statistics
```

## Features

- **Web Push Notifications**: Send push notifications to browser clients
- **Dashboard**: Real-time monitoring with:
  - Client statistics
  - Geographic distribution with interactive map
  - Top browsers and operating systems
  - Complete client list
- **Auto-detection**: Automatically detects client OS, browser, and location based off useragent and source ip
- **SQLite Database**: Persistent storage for subscriptions and statistics

## Setup

### 1. Install Go

Ensure Go 1.21+ is installed on your system.

### 2. Generate VAPID Keys

Before running the server, you need to generate VAPID keys:

1. Visit https://www.attheminute.com/vapid-key-generator
2. Click "Generate Keys"
3. Create a `data` folder in the project directory
4. Save the **Public Key** to `data/vapid_public.txt`
5. Save the **Private Key** to `data/vapid_private.txt`

Example:
```bash
mkdir data
# Copy your public key into data/vapid_public.txt
# Copy your private key into data/vapid_private.txt
```

### 3. Install Dependencies

```bash
go mod download
```

## Running the Server

```bash
go run .
```


The server will start on `http://localhost:10040`

**Access Points:**
- Dashboard: `http://localhost:10040`
- Network access: `http://YOUR_LOCAL_IP:10040` (shown in terminal output)

## Using the Dashboard

1. **Enable Notifications**: Click the "Enable Notifications" button in the header
2. **View Statistics**: See real-time client counts, geographic distribution, and browser/OS breakdown
3. **Send Notifications**: Use the "Send Notification" tab to broadcast messages to all subscribers
4. **Monitor Clients**: View detailed client list with IP, location, OS, and browser information

## Sending Notifications

### Via Web Dashboard
1. Navigate to the "Send Notification" tab in the sidebar
2. Fill in the notification title and message
3. Optionally add an icon URL
4. Click "Send to All Subscribers"

### Via API
Send a POST request to `/send-broadcast`:
```bash
curl -X POST http://localhost:10040/send-broadcast \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello","body":"Test notification","icon":""}'
```

## Configuration

- **Port**: Default is 10040, change in `main.go` if needed
- **Database**: SQLite database stored at `data/webpush.db`
- **VAPID Keys**: Must be manually placed in `data/` folder (see Setup section)

## Dependencies

- `github.com/SherClockHolmes/webpush-go` - Web Push protocol implementation
- `modernc.org/sqlite` - Pure Go SQLite implementation
- Leaflet.js - Interactive maps (loaded via CDN)

