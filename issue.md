# Issue: Implement Email Messaging Feature

## Objective
Membuat sistem email messaging yang modular dan reusable untuk berbagai kebutuhan aplikasi, seperti pengiriman notifikasi, verifikasi registrasi user, dan forgot password (reset password).

## Context
Saat ini sistem belum memiliki fitur untuk mengirimkan email. Fitur ini sangat penting untuk mendukung alur autentikasi (verifikasi email dan reset password) serta untuk mengirimkan notifikasi sistem ke pengguna.

## Requirements

### Functional Requirements
1. **Email Service**: Membuat service khusus yang menangani pengiriman email.
2. **Template Engine**: Mendukung HTML email templates sehingga email terlihat profesional dan dinamis (misalnya menggunakan `html/template`).
3. **Use Cases**:
   - Kirim Email Verifikasi Registrasi (berisi link/OTP).
   - Kirim Email Lupa Password (berisi link reset password).
   - Kirim Email Notifikasi Umum (misal: login dari device baru, perubahan password berhasil).

### Non-Functional Requirements
1. **Asynchronous Processing (Opsional/Direkomendasikan)**: Pengiriman email sebaiknya dilakukan secara asynchronous (misal menggunakan goroutine atau message queue) agar tidak memblokir HTTP response.
2. **Configurable**: Kredensial SMTP (Host, Port, Username, Password, Sender Name) harus diatur melalui `config.json` dan bisa di-load via environment variables atau config struct.
3. **Error Handling**: Log jika pengiriman email gagal tanpa membuat aplikasi crash.

## Technical Details / Proposed Implementation

### 1. Update Configuration
Tambahkan blok konfigurasi untuk SMTP pada `config.json` dan perbarui model config (`internal/config/config.go` atau sejenisnya):

```json
"smtp": {
  "host": "smtp.example.com",
  "port": 587,
  "username": "your_email@example.com",
  "password": "your_password",
  "sender_name": "GoFiber App",
  "sender_email": "noreply@example.com"
}
```

### 2. Dependencies
Gunakan library standar `net/smtp` yang sudah cukup mumpuni, atau bisa menggunakan third-party library seperti:
- `gopkg.in/gomail.v2` (direkomendasikan untuk handling attachment dan HTML dengan lebih mudah).
- `github.com/jordan-wright/email`

*Jika menggunakan gomail:*
`go get gopkg.in/gomail.v2`

### 3. Folder & File Structure
Buat struktur direktori untuk service dan template email:
```text
internal/
├── service/
│   ├── email_service.go      # Interface dan implementasi pengiriman email
│   └── email_service_test.go # Unit test untuk email service
templates/
└── email/
    ├── verification.html     # Template HTML verifikasi
    ├── reset_password.html   # Template HTML forgot password
    └── notification.html     # Template HTML notifikasi umum
```

### 4. Interface Design (`email_service.go`)
Buat interface agar mudah di-mock saat unit testing:
```go
package service

type EmailService interface {
	SendVerificationEmail(to, name, verificationLink string) error
	SendResetPasswordEmail(to, name, resetLink string) error
	SendNotificationEmail(to, subject, message string) error
}
```

## Task Checklist (Untuk Programmer / AI)

- [ ] **Task 1: Setup Config & Library**
  - Tambahkan struktur `SMTPConfig` pada konfigurasi aplikasi.
  - Tambahkan block `"smtp"` pada `config.json`.
  - Install dependensi (misal: `gomail.v2`) jika tidak memakai `net/smtp` bawaan.
- [ ] **Task 2: Buat HTML Templates**
  - Buat folder `templates/email/`.
  - Buat file template `verification.html`, `reset_password.html`, dan `notification.html`.
- [ ] **Task 3: Implementasi Email Service**
  - Buat file `internal/service/email_service.go`.
  - Implementasikan interface `EmailService`.
  - Buat fungsi helper untuk parsing HTML template dan mengganti variabel dinamis.
- [ ] **Task 4: Integrasi dengan Use Case (Authentication)**
  - Inject `EmailService` ke `AuthService` atau controller yang relevan.
  - Panggil `SendVerificationEmail` saat registrasi sukses.
  - Panggil `SendResetPasswordEmail` saat request forgot password valid.
- [ ] **Task 5: Testing & Validasi**
  - Tulis unit test sederhana untuk memastikan email builder / template parser berjalan lancar.
  - Tes pengiriman email (bisa menggunakan layanan seperti Mailtrap untuk development).

## Notes
- Harap menggunakan Mailtrap (`mailtrap.io`) atau layanan SMTP mock serupa untuk environment development agar tidak mengirimkan email asli secara tidak sengaja.
- Pastikan password SMTP tidak pernah di-commit ke dalam repository (gunakan environment variable di production).
