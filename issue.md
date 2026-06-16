# Implementasi Standalone RBAC Library untuk Fiber Microservices

## Latar Belakang
Sistem ERP dan aplikasi berskala *Enterprise* membutuhkan fitur **Role-Based Access Control (RBAC)** yang *granular* (tidak hanya sekadar level Admin/User, melainkan pengecekan hingga tingkat *action* seperti `sales.create`, `sales.approve`).

Untuk menjaga agar *boilerplate* / *microservices* tetap bersih, fitur RBAC ini akan diekstrak menjadi **dependency terpisah** (sebagai *Go Module*) sehingga bisa di-*consume* (di `go get`) oleh *project* manapun yang menggunakan *framework* Go-Fiber.

## Tujuan / Scope Pekerjaan
Membuat *Go Module* independen (misal: `github.com/fatihrizqon/gofiber-rbac`) yang mengekspos **Fiber Middleware** untuk melakukan verifikasi *permission* dari sebuah *request*. Library ini **TIDAK** boleh terkunci pada satu jenis *Database* (Postgres/MySQL) sehingga ia harus menerima abstraksi berupa `Interface` untuk pengambilan data *permission*.

---

## Spesifikasi Teknis (Untuk Dieksekusi oleh Programmer / AI Model)

### 1. Inisialisasi Project
- Buat folder terpisah dari *boilerplate* ini (misal: `gofiber-rbac`).
- Jalankan `go mod init github.com/fatihrizqon/gofiber-rbac`.
- Tambahkan dependensi `github.com/gofiber/fiber/v2`.

### 2. Desain Interface (Contract)
Library ini tidak peduli dari mana data *permission* berasal (entah itu DB, Redis, atau file statis). Oleh karena itu, siapkan sebuah *Interface* `PermissionStore` yang nantinya **harus diimplementasikan** oleh *Consumer* (Aplikasi utama).

```go
package rbac

// PermissionStore adalah interface yang wajib di-inject oleh aplikasi utama.
type PermissionStore interface {
	// GetUserPermissions mengembalikan daftar permission (string) yang dimiliki oleh user ID tertentu.
	// Contoh return: ["sales.read", "sales.create"]
	GetUserPermissions(userId string) ([]string, error)
}
```

### 3. Konfigurasi Library
Buat *struct* konfigurasi agar implementasinya fleksibel. Terkadang *Consumer* menyimpan ID User di *Context* (Locals) dengan *key* yang berbeda (misal: `"user_id"` atau JWT token object).

```go
package rbac

import "github.com/gofiber/fiber/v2"

type Config struct {
	Store PermissionStore
	
	// UserLookup mendefinisikan cara mengambil UserID dari Fiber Context (setelah melewati JWT Middleware).
	// Default: func(c *fiber.Ctx) string { return c.Locals("user_id").(string) }
	UserLookup func(c *fiber.Ctx) string 

	// UnauthorizedHandler adalah custom response jika user tidak memiliki akses
	UnauthorizedHandler fiber.Handler
}
```

### 4. Implementasi Middleware Utama
Buat *constructor* dan *middleware generator*.

```go
package rbac

import "github.com/gofiber/fiber/v2"

type RBAC struct {
	config Config
}

func New(cfg Config) *RBAC {
    // Set default UserLookup dan UnauthorizedHandler jika tidak diisi oleh consumer
	// ...
	return &RBAC{config: cfg}
}

// Require adalah middleware yang akan dipasang pada Route Fiber
func (r *RBAC) Require(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil userID dari Context menggunakan config.UserLookup(c)
		// 2. Jika userID kosong, panggil UnauthorizedHandler
		
		// 3. Ambil list permissions dari config.Store.GetUserPermissions(userID)
		// 4. Lakukan pengecekan: apakah "permission" ada di dalam list tersebut?
		
		// 5. Jika ADA: return c.Next()
		// 6. Jika TIDAK ADA: panggil UnauthorizedHandler
		return c.Next() // ganti dengan logika yang benar
	}
}
```

### 5. Expected Usage (Contoh Penggunaan oleh Consumer)
Library yang dibuat harus bisa digunakan semudah ini di aplikasi utama:

```go
// Di dalam aplikasi utama (contoh: gofiber-microservice)

// 1. Buat Adapter
type MyRBACStore struct {
    db *gorm.DB // atau Redis
}
func (s *MyRBACStore) GetUserPermissions(userId string) ([]string, error) {
    // Query ke DB untuk mengambil permission...
    return []string{"sales.create"}, nil
}

// 2. Setup RBAC
rbacStore := &MyRBACStore{db: database}
rbacEngine := rbac.New(rbac.Config{
    Store: rbacStore,
    UserLookup: func(c *fiber.Ctx) string {
        // Asumsi JWT payload ditaruh di Locals
        user := c.Locals("user").(*jwt.Token)
        claims := user.Claims.(jwt.MapClaims)
        return claims["id"].(string)
    },
})

// 3. Pasang di Route
app.Post("/sales", authMiddleware, rbacEngine.Require("sales.create"), salesHandler.Create)
```

### 6. Testing (Wajib)
Buat file `rbac_test.go`. Tuliskan *Unit Test* dengan `httptest` Fiber:
- Uji saat User ID kosong.
- Uji saat User tidak punya *permission* (harus 403 Forbidden).
- Uji saat User punya *permission* (harus meneruskan ke *handler* sukses).

---

## Kriteria Penerimaan (Acceptance Criteria)
- [ ] Kode terisolasi dalam satu module tanpa *logic business* dari aplikasi spesifik.
- [ ] Interface `PermissionStore` didefinisikan dengan benar.
- [ ] Middleware berfungsi melakukan *block* (HTTP 403) jika *permission* tidak cocok.
- [ ] Terdapat `README.md` singkat mengenai cara *install* dan penggunaan library.
- [ ] 100% *Unit test coverage* untuk fungsi `Require()`.
