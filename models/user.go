package models

import "time"

// User merepresentasikan data user untuk login
type User struct {
	ID        int       `json:"id"`
	RfidUID   *string   `json:"rfid_uid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Nama      string    `json:"nama_lengkap"`
	Role      string    `json:"role"`   // admin atau karyawan
	Status    string    `json:"status"` // aktif atau nonaktif
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
