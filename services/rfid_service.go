package services

import (
	"log"

	"github.com/yourusername/presensi-agent/config"
	"github.com/yourusername/presensi-agent/models"
)

// FindUserByRFID mencari user berdasarkan rfid_uid
func FindUserByRFID(rfidUID string) (*models.User, error) {
	query := "SELECT id, rfid_uid, username, email, nama_lengkap, role, status FROM users WHERE rfid_uid = ? AND status = 'aktif'"

	row := config.DB.QueryRow(query, rfidUID)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.RfidUID, &user.Username, &user.Email, &user.Nama, &user.Role, &user.Status)

	if err != nil {
		log.Printf("User dengan RFID %s tidak ditemukan: %v", rfidUID, err)
		return nil, err
	}

	return user, nil
}

// ProcessRFIDData memproses data RFID yang diterima dari WebSocket
func ProcessRFIDData(rfidData *models.RFIDData) (*models.User, error) {
	log.Printf("Processing RFID data: uid=%s, timestamp=%s", rfidData.UID, rfidData.Timestamp)

	// Cari user dengan rfid_uid yang cocok
	user, err := FindUserByRFID(rfidData.UID)
	if err != nil {
		log.Printf("❌ RFID tidak cocok: %s - Data tidak akan direcord", rfidData.UID)
		return nil, err
	}

	log.Printf("✅ RFID cocok! User ditemukan: %s (ID: %d)", user.Nama, user.ID)
	return user, nil
}
