# Presensi Agent - Aplikasi Presensi dengan RFID WebSocket

Aplikasi Go untuk menangkap data RFID dari WebSocket, memvalidasi dengan database, dan merekam presensi otomatis.

## Fitur

✅ WebSocket client untuk menangkap data RFID real-time
✅ Validasi RFID dengan database (rfid_uid di table users)
✅ Auto record presensi ke database
✅ Mendeteksi jam masuk dan jam keluar
✅ Error handling yang baik
✅ Logging yang informatif

## Struktur Project

```
presensi-agent/
├── main.go                 # Entry point aplikasi
├── go.mod                  # Go module definition
├── .env.example            # Contoh environment variables
├── config/
│   └── database.go         # Konfigurasi koneksi MySQL
├── models/
│   ├── user.go            # Struct User
│   └── presensi.go        # Struct Presensi & RFIDData
├── services/
│   ├── rfid_service.go    # Service untuk RFID matching
│   └── presensi_service.go # Service untuk recording presensi
└── README.md
```

## Installation

1. Clone repository atau copy project files
2. Install dependencies:
   ```bash
   go mod download
   ```

3. Copy `.env.example` menjadi `.env` dan edit sesuai konfigurasi database:
   ```bash
   cp .env.example .env
   ```

4. Pastikan database sudah dibuat dengan script `presensi_db.sql`

## Configuration

Edit file `.env`:

```env
# Database
DB_HOST=103.175.218.246
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=presensi_db

# WebSocket
WS_URL=ws://localhost:8080
WS_PATH=/ws/rfid
```

## Cara Menjalankan

```bash
go run main.go
```

Atau build executable:
```bash
go build -o presensi-agent.exe
./presensi-agent.exe
```

## Alur Kerja

1. **Connect ke WebSocket**: Aplikasi terhubung ke server WebSocket yang mengirim data RFID
2. **Terima Data**: Data RFID diterima dalam format JSON:
   ```json
   {
     "uid": "F1:1B:C3:01",
     "length": 4,
     "timestamp": "09:09:46",
     "unix_time": 1765764586
   }
   ```

3. **Validasi RFID**: UID dari RFID dicocokkan dengan `rfid_uid` di table `users`
   - Jika **cocok** → Lanjut ke step 4
   - Jika **tidak cocok** → **Tidak record presensi** (kosong)

4. **Record Presensi**: 
   - Jika belum ada presensi hari ini → Insert baru dengan `jam_masuk`
   - Jika sudah ada presensi hari ini → Update dengan `jam_keluar`

5. **Logging**: Setiap event dicatat dengan log yang detail

## Contoh Output

```
2025-12-15 14:30:45 ✅ WebSocket connected successfully!
2025-12-15 14:30:50 📨 Message received: {"uid":"F1:1B:C3:01","length":4,"timestamp":"09:09:46","unix_time":1765764586}
2025-12-15 14:30:50 Processing RFID data: uid=F1:1B:C3:01, timestamp=09:09:46
2025-12-15 14:30:50 ✅ RFID cocok! User ditemukan: Budi Santoso (ID: 1)
2025-12-15 14:30:50 ✅ Presensi recorded (jam_masuk) untuk user 1
2025-12-15 14:30:50 📤 Response sent: {"status":"success","message":"Presensi recorded",...}
```

## Database Schema

### Table: users
- id (PRIMARY KEY)
- rfid_uid (untuk matching RFID)
- username
- email
- password
- nama_lengkap
- role (admin/karyawan)
- status (aktif/nonaktif)

### Table: presensi
- id (PRIMARY KEY)
- user_id (FOREIGN KEY)
- tanggal
- jam_masuk
- jam_keluar
- status (hadir/sakit/izin/libur/tidak_hadir)
- keterangan

## Error Handling

- **RFID tidak cocok**: Logged sebagai error, tidak akan record presensi
- **Database error**: Logged dan tidak merekam presensi
- **JSON parsing error**: Logged dan skip message
- **WebSocket disconnect**: Akan disconnect dan exit dengan error

## Dependencies

- `github.com/gorilla/websocket` - WebSocket client
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/joho/godotenv` - Load environment variables

## Lisensi

MIT
