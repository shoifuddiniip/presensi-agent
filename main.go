package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/yourusername/presensi-agent/config"
	"github.com/yourusername/presensi-agent/models"
	"github.com/yourusername/presensi-agent/services"
)

func init() {
	// Load environment variables dari .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  File .env tidak ditemukan, menggunakan default env variables")
	}
}

func main() {
	// Initialize database connection
	err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDatabase()

	// Get WebSocket URL dari environment
	wsURL := os.Getenv("WS_URL")
	if wsURL == "" {
		wsURL = "ws://localhost:8080"
	}

	wsPath := os.Getenv("WS_PATH")
	if wsPath == "" {
		wsPath = "/ws/rfid"
	}

	fullURL := wsURL + wsPath
	log.Printf("Connecting to WebSocket: %s", fullURL)

	// Connect to WebSocket server
	header := make(url.Values)
	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}

	conn, _, err := dialer.Dial(fullURL, header)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	log.Println("✅ WebSocket connected successfully!")

	// Goroutine untuk keep-alive/heartbeat
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("Failed to send heartbeat: %v", err)
				return
			}
		}
	}()

	// Listen untuk messages dari WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("📨 Message received: %s", string(message))

		// Parse RFID data dari JSON
		var rfidData models.RFIDData
		err = json.Unmarshal(message, &rfidData)
		if err != nil {
			log.Printf("❌ Error parsing JSON: %v", err)
			continue
		}

		// Process RFID data - cari user berdasarkan RFID UID
		user, err := services.ProcessRFIDData(&rfidData)
		if err != nil {
			// RFID tidak cocok dengan user manapun - jangan record presensi
			continue
		}

		// RFID cocok - record presensi untuk user
		presensi, err := services.RecordPresensi(user.ID, &rfidData)
		if err != nil {
			log.Printf("❌ Error recording presensi: %v", err)
			continue
		}

		// Send response
		response := map[string]interface{}{
			"status":   "success",
			"message":  "Presensi recorded",
			"user_id":  user.ID,
			"username": user.Username,
			"nama":     user.Nama,
			"presensi": presensi,
		}

		responseJSON, _ := json.Marshal(response)
		log.Printf("📤 Response sent: %s", string(responseJSON))
	}

	log.Println("WebSocket connection closed")
}
