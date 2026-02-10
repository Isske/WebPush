# WebPush Architecture

## Component Overview

### Main Application (main.go)
- Entry point
- Initializes components
- Sets up HTTP routes
- Starts background services

### Handlers (handlers/)
Organized by functionality:

#### vapid.go
- VAPID key generation and management
- Stores keys in `data/vapid_*.txt`
- Exposes public key via `/vapid-public-key` endpoint

#### subscription.go
- Handles new subscriber registrations
- Collects client metadata (IP, geo, browser, OS)
- Stores subscriptions in `data/subscriptions.json`
- Endpoint: `/subscribe`

#### notification.go
- Sends push notifications via API
- Terminal-based notification sender
- Tracks push counts in `data/push_count.json`
- Endpoint: `/send-notification`

#### dashboard.go
- Serves dashboard interface
- Provides statistics API
- Aggregates subscription data
- Endpoints: `/` and `/api/stats`

### Models (models/)
Defines shared data structures:
- `Subscription` - Client subscription with metadata
- `NotificationPayload` - Push notification content
- `SendRequest` - API request format
- `DashboardStats` - Dashboard statistics

### Utilities (utils/)

#### useragent.go
- Parses User-Agent strings
- Extracts OS name and version
- Extracts browser name and version
- Handles platform-specific version parsing

#### geoip.go
- Performs GeoIP lookups using ip-api.com
- Returns country code for IP addresses
- Non-blocking, best-effort approach

### Static Assets (static/)

#### index.html
- Dashboard web interface
- Dark theme with Lora font
- Leaflet.js choropleth map
- Real-time statistics display

#### sw.js
- Service worker for push notifications
- Handles notification display
- Manages background notifications

## Data Flow

### Subscription Flow
```
Browser → /subscribe → HandleSubscribe
                       ↓
                   Parse UserAgent
                       ↓
                   Lookup GeoIP
                       ↓
                Save to subscriptions.json
```

### Notification Flow
```
API/Terminal → /send-notification → SendNotificationHandler
                                    ↓
                            Create webpush.Subscription
                                    ↓
                            Marshal payload to JSON
                                    ↓
                            Send via webpush library
                                    ↓
                            Increment push count
```

### Dashboard Flow
```
Browser → /api/stats → GetDashboardStatsHandler
                       ↓
                   Load subscriptions
                       ↓
                   Aggregate by country
                       ↓
                   Load push count
                       ↓
                   Return JSON response
```

## Port Configuration

Default: `10040`

All services run on a single port:
- Dashboard UI: `http://localhost:10040/`
- API endpoints: `http://localhost:10040/api/*`
- Push endpoints: `http://localhost:10040/subscribe`, etc.

## File Locations

### Generated at Runtime
- `data/vapid_public.txt` - VAPID public key
- `data/vapid_private.txt` - VAPID private key
- `data/subscriptions.json` - All subscriptions
- `data/push_count.json` - Push statistics

### Static Assets
- `static/index.html` - Dashboard UI
- `static/sw.js` - Service worker

## Security Considerations

- VAPID keys are auto-generated and stored locally
- No authentication on endpoints (add if needed)
- GeoIP lookups use external API (best-effort)
- Subscription data stored in plain JSON


