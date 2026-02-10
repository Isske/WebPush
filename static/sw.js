// Service Worker for handling push notifications
self.addEventListener('push', event => {
    console.log('[sw.js] Push event received:', event);
    if (event.data) {
        try {
            const raw = event.data.text();
            console.log('[sw.js] Push event raw data:', raw);
            let data = {};
            try {
                data = JSON.parse(raw);
            } catch (e) {
                console.error('[sw.js] Error parsing JSON:', e);
            }
            const options = {
                body: data.body || raw,
                icon: data.icon || '/static/icon.png',
                badge: data.badge || '/static/badge.png',
                tag: data.tag || 'notification',
                vibrate: data.vibrate || [200, 100, 200],
                requireInteraction: false,
            };
            event.waitUntil(
                self.registration.showNotification(data.title || 'Notification', options)
                    .catch(err => {
                        console.error('[sw.js] showNotification error:', err);
                        return self.registration.showNotification('Notification Fallback', {
                            body: '[sw.js] showNotification error: ' + err.message,
                        });
                    })
            );
        } catch (error) {
            console.error('[sw.js] Error handling push event:', error);
            event.waitUntil(
                self.registration.showNotification('Notification', {
                    body: '[sw.js] Error: ' + error.message,
                })
            );
        }
    } else {
        console.warn('[sw.js] Push event with no data!');
        event.waitUntil(
            self.registration.showNotification('Notification', {
                body: '[sw.js] Push event with no data',
            })
        );
    }
});

// Handle notification clicks
self.addEventListener('notificationclick', event => {
    console.log('Notification clicked:', event);
    event.notification.close();
    
    event.waitUntil(
        clients.matchAll({ type: 'window', includeUncontrolled: true }).then(clientList => {
            // Check if there's already a window open
            for (let i = 0; i < clientList.length; i++) {
                const client = clientList[i];
                if (client.url === '/' && 'focus' in client) {
                    return client.focus();
                }
            }
            // If not, open a new window
            if (clients.openWindow) {
                return clients.openWindow('/');
            }
        })
    );
});

// Handle notification close
self.addEventListener('notificationclose', event => {
    console.log('Notification closed:', event);
});
