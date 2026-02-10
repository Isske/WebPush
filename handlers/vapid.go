package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	PublicKeyFile  = "data/vapid_public.txt"
	PrivateKeyFile = "data/vapid_private.txt"
)

var (
	VapidPublicKey  string
	VapidPrivateKey string
)

// InitVAPIDKeys loads VAPID keys from files
func InitVAPIDKeys() error {
	// Read VAPID public key
	pubKeyData, err := os.ReadFile(PublicKeyFile)
	if err != nil {
		log.Fatalf("Failed to read VAPID public key from %s: %v\n", PublicKeyFile, err)
		log.Println("Please generate VAPID keys at https://www.attheminute.com/vapid-key-generator")
		log.Println("and place them in data/vapid_public.txt and data/vapid_private.txt")
		return err
	}
	VapidPublicKey = strings.TrimSpace(string(pubKeyData))

	// Read VAPID private key
	privKeyData, err := os.ReadFile(PrivateKeyFile)
	if err != nil {
		log.Fatalf("Failed to read VAPID private key from %s: %v\n", PrivateKeyFile, err)
		log.Println("Please generate VAPID keys at https://www.attheminute.com/vapid-key-generator")
		log.Println("and place them in data/vapid_public.txt and data/vapid_private.txt")
		return err
	}
	VapidPrivateKey = strings.TrimSpace(string(privKeyData))

	// Validate keys are not empty
	if VapidPublicKey == "" || VapidPrivateKey == "" {
		log.Fatal("VAPID keys are empty. Please generate keys at https://www.attheminute.com/vapid-key-generator")
	}

	log.Println("✓ VAPID keys loaded successfully")
	log.Printf("Public Key: %s", VapidPublicKey)
	return nil
}

// GetVAPIDPublicKeyHandler returns the VAPID public key
func GetVAPIDPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(VapidPublicKey))
}
