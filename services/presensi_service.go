package services

import (
	"database/sql"
	"log"
	"time"

	"github.com/yourusername/presensi-agent/config"
	"github.com/yourusername/presensi-agent/models"
)

// RecordPresensi merekam presensi user ke database
func RecordPresensi(userID int, rfidData *models.RFIDData) (*models.Presensi, error) {
	// Ambil tanggal hari ini
	today := time.Now().Format("2006-01-02")

	// Cek apakah sudah ada presensi untuk user hari ini
	checkQuery := "SELECT id FROM presensi WHERE user_id = ? AND tanggal = ?"
	existingRow := config.DB.QueryRow(checkQuery, userID, today)

	var existingID int
	err := existingRow.Scan(&existingID)

	if err == nil {
		// Presensi sudah ada, update dengan jam_keluar
		log.Printf("Presensi untuk user %d sudah ada hari ini. Update jam_keluar...", userID)

		currentTime := time.Now().Format("15:04:05")
		updateQuery := "UPDATE presensi SET jam_keluar = ?, updated_at = NOW() WHERE id = ?"

		result, err := config.DB.Exec(updateQuery, currentTime, existingID)
		if err != nil {
			log.Printf("Error updating presensi: %v", err)
			return nil, err
		}

		rows, _ := result.RowsAffected()
		log.Printf("✅ Presensi updated (jam_keluar): %d row affected", rows)

		return &models.Presensi{
			ID:     existingID,
			UserID: userID,
		}, nil

	} else if err != sql.ErrNoRows {
		// Error selain data tidak ditemukan
		log.Printf("Error checking existing presensi: %v", err)
		return nil, err
	}

	// Presensi belum ada, buat yang baru dengan jam_masuk
	currentTime := time.Now().Format("15:04:05")
	insertQuery := `
		INSERT INTO presensi (user_id, tanggal, jam_masuk, status, created_at, updated_at)
		VALUES (?, ?, ?, 'hadir', NOW(), NOW())
	`

	result, err := config.DB.Exec(insertQuery, userID, today, currentTime)
	if err != nil {
		log.Printf("Error recording presensi: %v", err)
		return nil, err
	}

	id, _ := result.LastInsertId()
	log.Printf("✅ Presensi recorded (jam_masuk) untuk user %d", userID)

	return &models.Presensi{
		ID:       int(id),
		UserID:   userID,
		Tanggal:  today,
		JamMasuk: &currentTime,
		Status:   "hadir",
	}, nil
}

// GetUserPresensiToday mendapatkan presensi user hari ini
func GetUserPresensiToday(userID int) (*models.Presensi, error) {
	today := time.Now().Format("2006-01-02")
	query := "SELECT id, user_id, tanggal, jam_masuk, jam_keluar, status FROM presensi WHERE user_id = ? AND tanggal = ?"

	row := config.DB.QueryRow(query, userID, today)

	presensi := &models.Presensi{}
	err := row.Scan(&presensi.ID, &presensi.UserID, &presensi.Tanggal, &presensi.JamMasuk, &presensi.JamKeluar, &presensi.Status)

	if err != nil {
		return nil, err
	}

	return presensi, nil
}
