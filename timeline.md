# Timeline & Execution Plan for Accounting PRD
**Target Executor:** Junior Programmer / Low-Model AI
**Goal:** Implement `accounting_prd_v1.md` securely and systematically without architectural violations.

> [!WARNING]
> **Aturan Emas untuk Eksekutor (Junior/AI):**
> 1. **Gunakan Boilerplate Pattern:** Wajib memisahkan `delivery` (HTTP), `service` (Logika Bisnis), dan `repository` (Database).
> 2. **Tenant Isolation:** Untuk semua tabel selain tabel sistem (seperti `users`, `roles`), parameter `company_id` **wajib** disuntikkan dari `c.Locals` ke `repository`.
> 3. **ACID Transactions:** Dilarang melakukan multiple insert/update tanpa fungsi *Database Transaction* (Tx).
> 4. **Review Wajib:** Jangan lanjut ke fase berikutnya sebelum kode di-review oleh Senior Engineer.

---

## 🏗️ Phase 1: Foundation & Authentication (Minggu 1)
**Fokus:** Membangun lapisan akses dan autentikasi global. Kesalahan di sini akan menyebabkan celah keamanan.

- [ ] **Task 1.1: Database Migration**
  - Buat file migrasi (DDL) untuk **semua** tabel yang ada di PRD (menggunakan alat migrasi seperti `golang-migrate` atau ekivalennya).
  - Tulis skema tabel persis seperti di PRD (perhatikan tipe data `UUID`, constraints, dan foreign keys).
  - *Checkpoint:* Senior melakukan review ERD & Skema DB.

- [ ] **Task 1.2: Skema RBAC & Seeding**
  - Buat struktur endpoint CRUD untuk `roles` dan `permissions`.
  - Buat script *seeding* (atau endpoint inisialisasi) untuk memasukkan "Default Role-Permission Matrix".

- [ ] **Task 1.3: Auth & User Module**
  - Implementasikan registrasi, login, refresh token, dan logout.
  - Implementasikan JWT dual-token ke dalam `AuthMiddleware`.
  - Pastikan pengamanan *bcrypt* untuk password dan penanganan sesi (revocation).

---

## 🏢 Phase 2: Multi-Tenancy & Master Data (Minggu 2)
**Fokus:** Membangun modul-modul CRUD sederhana dengan pengamanan tenant (`company_id`).

- [ ] **Task 2.1: Company & Members**
  - Buat fitur pembuatan perusahaan. Otomatis tambahkan pembuatnya sebagai `owner` di `company_members` menggunakan `role_id` dari DB.
  - Selesaikan `CompanyMiddleware` dan `PermissionMiddleware` untuk memvalidasi akses *role* berdasarkan data di *database*.

- [ ] **Task 2.2: Chart of Accounts (COA)**
  - Implementasikan CRUD untuk COA Group, Subgroup, dan Accounts. 
  - Pastikan setiap data diisolasi menggunakan `company_id`.

- [ ] **Task 2.3: Fiscal Year & Master Data Pendukung**
  - Implementasikan modul Fiscal Year (Tahun Fiskal) dan Periods.
  - Implementasikan modul Customers (Pelanggan) dan Vendors (Pemasok).
  - *Checkpoint:* Lakukan *integration test* untuk memastikan pengguna dari Perusahaan A tidak bisa membaca COA/Vendor dari Perusahaan B.

---

## ⚙️ Phase 3: The Core Accounting Engine (Minggu 3) - ⚠️ CRITICAL
**Fokus:** Menangani *Double-Entry Accounting*. Ini adalah jantung aplikasi, jangan gunakan cara manipulasi data sederhana.

- [ ] **Task 3.1: Transaction Pattern (Boilerplate Update)**
  - Junior/AI harus diberikan 1 contoh konkrit fungsi `Repository` yang menangani `Tx` (Database Transaction). 

- [ ] **Task 3.2: Journal Entries & Journal Lines**
  - Buat modul `journal`.
  - **Validasi Mutlak:** Total `debit` di `journal_lines` **harus persis sama** dengan total `credit`. Jika tidak, kembalikan HTTP 422.
  - **Presisi Mutlak:** Gunakan library `shopspring/decimal` (atau sejenisnya) untuk struktur model Go, bukan `float64`.
  - Simpan `journal_entries` dan `journal_lines` dalam 1 transaksi DB. Jika satu gagal, Rollback semua.
  - *Checkpoint:* Senior melakukan uji coba *stress test* dan pengecekan *rounding error* pada nominal triliunan rupiah.

---

## 🛒 Phase 4: Procurement & Sales Flow (Minggu 4)
**Fokus:** Otomatisasi Jurnal dari alur transaksi jual-beli. Menghindari *spaghetti code*.

- [ ] **Task 4.1: Inter-Module Service Communication**
  - Siapkan *Dependency Injection* atau pola berbasis *Event/Outbox* agar `SalesService` dapat memanggil `JournalService`.

- [ ] **Task 4.2: Procurement (Purchasing)**
  - Implementasi CRUD Purchase Orders.
  - Implementasi konversi PO menjadi Bill (Tagihan).
  - Implementasi Bill Payment. **Aturan:** Saat Bill dibayar, sistem harus memanggil `JournalService` untuk mencatat Kas/Bank (Kredit) dan Hutang/Account Payable (Debit) secara terprogram.

- [ ] **Task 4.3: Sales (Invoicing)**
  - Implementasi CRUD Quotation dan Invoice.
  - Saat Invoice dilunasi (Payment), sistem memanggil `JournalService` untuk mencatat Kas/Bank (Debit) dan Piutang (Kredit).

---

## 📊 Phase 5: Reporting & Aggregation (Minggu 5)
**Fokus:** Membaca data secara optimal tanpa menyebabkan server *crash*. Dilarang menggunakan *N+1 Queries*.

- [ ] **Task 5.1: Database Views / Complex Queries**
  - Senior menyediakan query SQL atau *Recursive CTE* yang efisien untuk membaca laporan. AI/Junior bertugas menerjemahkannya ke lapisan `Repository`.

- [ ] **Task 5.2: Financial Reports**
  - Buat endpoint `Trial Balance`, `Profit & Loss`, dan `Balance Sheet` (`read-only`).
  
- [ ] **Task 5.3: Dashboard Aggregation**
  - Buat endpoint untuk *summary* ringkas (misal: Total Hutang Jatuh Tempo, Total Pendapatan Bulan Ini).

---

## 🏁 Phase 6: Final Review & UAT
- Pembersihan *codebase*.
- Memastikan *Swagger Docs* (`swag init`) sukses digenerate dan mewakili semua endpoint.
- Validasi fungsional menyeluruh (User Acceptance Testing).
