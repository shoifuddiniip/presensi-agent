package models

import "time"

// Presensi merepresentasikan data kehadiran
type Presensi struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Tanggal    string    `json:"tanggal"`
	JamMasuk   *string   `json:"jam_masuk"`
	JamKeluar  *string   `json:"jam_keluar"`
	Status     string    `json:"status"` // hadir, sakit, izin, libur, tidak_hadir
	Keterangan *string   `json:"keterangan"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// RFIDData merepresentasikan data yang diterima dari WebSocket RFID reader
type RFIDData struct {
	UID       string `json:"uid"`
	Length    int    `json:"length"`
	Timestamp string `json:"timestamp"`
	UnixTime  int64  `json:"unix_time"`
}
