package models

// Subscription represents a web push subscription with client metadata
type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
	IP              string `json:"ip,omitempty"`
	Nation          string `json:"nation,omitempty"`
	OS              string `json:"os,omitempty"`
	OSVersion       string `json:"os_version,omitempty"`
	Browser         string `json:"browser,omitempty"`
	BrowserVersion  string `json:"browser_version,omitempty"`
	Platform        string `json:"platform,omitempty"`
	PlatformVersion string `json:"platform_version,omitempty"`
}

// NotificationPayload defines the structure of a push notification
type NotificationPayload struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	Icon    string `json:"icon"`
	Badge   string `json:"badge"`
	Tag     string `json:"tag"`
	Vibrate []int  `json:"vibrate"`
}

// SendRequest is the request format for sending a notification
type SendRequest struct {
	Subscription Subscription `json:"subscription"`
	Title        string       `json:"title"`
	Body         string       `json:"body"`
	Icon         string       `json:"icon"`
	Badge        string       `json:"badge"`
}

// DashboardStats aggregates statistics for the dashboard view
type DashboardStats struct {
	TotalClients  int            `json:"total_clients"`
	OnlineClients int            `json:"online_clients"`
	TotalPushes   int            `json:"total_pushes"`
	Countries     map[string]int `json:"countries"`
	Subscriptions []Subscription `json:"subscriptions"`
}
