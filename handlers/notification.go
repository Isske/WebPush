package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	"webpush/database"
	"webpush/models"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// SendNotificationHandler handles API requests to send push notifications
func SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SendRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Failed to decode request: %v\n", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("[Push] Sending notification to endpoint: %s", req.Subscription.Endpoint)

	// Create the push subscription
	s := &webpush.Subscription{
		Endpoint: req.Subscription.Endpoint,
		Keys: webpush.Keys{
			P256dh: req.Subscription.Keys.P256dh,
			Auth:   req.Subscription.Keys.Auth,
		},
	}

	// Create notification payload
	payload := models.NotificationPayload{
		Title:   req.Title,
		Body:    req.Body,
		Icon:    req.Icon,
		Badge:   req.Badge,
		Vibrate: []int{200, 100, 200},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v\n", err)
		http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
		return
	}

	// Send the notification
	resp, err := webpush.SendNotification(payloadJSON, s, &webpush.Options{
		Subscriber:      "mailto:example@example.com",
		VAPIDPublicKey:  VapidPublicKey,
		VAPIDPrivateKey: VapidPrivateKey,
		TTL:             30,
	})

	if err != nil {
		log.Printf("[Push] Error sending notification: %v", err)
		http.Error(w, "Failed to send notification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[Push] Notification response status: %d", resp.StatusCode)

	if resp.StatusCode == 404 || resp.StatusCode == 410 {
		log.Printf("[Push] Subscription is no longer valid (status %d). Removing from disk.", resp.StatusCode)
		RemoveLatestSubscription()
		LatestSubscription = nil
		http.Error(w, "Subscription is no longer valid", http.StatusGone)
		return
	}

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[Push] Error response body: %s", string(bodyBytes))
		http.Error(w, "Failed to send notification: "+string(bodyBytes), resp.StatusCode)
		return
	}

	log.Printf("[Push] Notification sent successfully")
	database.IncrementPushCount(1)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

// StartAutoSender starts a background goroutine that sends notifications periodically
func StartAutoSender() {
	go func() {
		for {
			time.Sleep(time.Minute)

			if LatestSubscription != nil {
				log.Println("[AutoPush] Sending scheduled notification...")
				s := &webpush.Subscription{
					Endpoint: LatestSubscription.Endpoint,
					Keys: webpush.Keys{
						P256dh: LatestSubscription.Keys.P256dh,
						Auth:   LatestSubscription.Keys.Auth,
					},
				}

				payload := models.NotificationPayload{
					Title:   "Scheduled Notification",
					Body:    "This is an automatic notification from the backend.",
					Icon:    "",
					Badge:   "",
					Vibrate: []int{200, 100, 200},
				}

				payloadJSON, _ := json.Marshal(payload)
				resp, err := webpush.SendNotification(payloadJSON, s, &webpush.Options{
					Subscriber:      "mailto:example@example.com",
					VAPIDPublicKey:  VapidPublicKey,
					VAPIDPrivateKey: VapidPrivateKey,
					TTL:             30,
				})

				if err != nil {
					log.Printf("[AutoPush] Error sending notification: %v\n", err)
				} else {
					log.Printf("[AutoPush] Notification response status: %d\n", resp.StatusCode)
				}
			}
		}
	}()
}

// SendBroadcastHandler sends a notification to all subscriptions
func SendBroadcastHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Icon  string `json:"icon"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Failed to decode broadcast request: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	subs := LoadSubscriptions()
	if len(subs) == 0 {
		log.Println("[Broadcast] No subscriptions found.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"sent":   0,
			"failed": 0,
		})
		return
	}

	payload := models.NotificationPayload{
		Title:   req.Title,
		Body:    req.Body,
		Icon:    req.Icon,
		Badge:   "",
		Vibrate: []int{200, 100, 200},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Broadcast] Error marshaling payload: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to marshal payload"})
		return
	}

	log.Printf("[Broadcast] Sending to %d subscriptions.", len(subs))
	var validSubs []models.Subscription
	sent := 0
	failed := 0

	for _, sub := range subs {
		s := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.Keys.P256dh,
				Auth:   sub.Keys.Auth,
			},
		}

		resp, err := webpush.SendNotification(payloadJSON, s, &webpush.Options{
			Subscriber:      "mailto:example@example.com",
			VAPIDPublicKey:  VapidPublicKey,
			VAPIDPrivateKey: VapidPrivateKey,
			TTL:             30,
		})

		if err != nil {
			log.Printf("[Broadcast] Error sending to %s: %v", sub.Endpoint, err)
			failed++
			continue
		}

		log.Printf("[Broadcast] Response status for %s: %d", sub.Endpoint, resp.StatusCode)

		if resp.StatusCode == 404 || resp.StatusCode == 410 {
			log.Printf("[Broadcast] Subscription %s is no longer valid (status %d). Removing.", sub.Endpoint, resp.StatusCode)
			failed++
			continue
		}

		if resp.StatusCode >= 400 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Printf("[Broadcast] Error response body for %s: %s", sub.Endpoint, string(bodyBytes))
			failed++
			continue
		}

		log.Printf("[Broadcast] Notification sent to: %s", sub.Endpoint)
		validSubs = append(validSubs, sub)
		sent++
	}

	database.IncrementPushCount(sent)
	SaveSubscriptions(validSubs)

	if len(validSubs) == 0 {
		LatestSubscription = nil
	} else {
		LatestSubscription = &validSubs[len(validSubs)-1]
	}

	log.Printf("[Broadcast] Sent: %d, Failed: %d", sent, failed)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sent":   sent,
		"failed": failed,
	})
}
