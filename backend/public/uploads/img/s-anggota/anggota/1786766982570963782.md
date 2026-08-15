# DPP Gradasi — Backend

Backend Go (Gin + GORM + MySQL) untuk aplikasi web DPP GRADASI.
Arsitektur modular per domain: `internal/module/<nama>/{dto,entity,handler,mapper,model,repository,routes,service}`.

## Persyaratan

- Go 1.21+
- MySQL 8+
- [goose](https://github.com/pressly/goose) untuk migrasi (opsional — bisa lewat Makefile)

## Setup

```bash
# 1. Salin konfigurasi
cp .env.example .env
# 2. Isi .env sesuai environment Anda (DB, JWT_SECRET, SMTP, dsb)

# 3. Buat database (sekali)
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS dpp_gradasi_jakarta CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 4. Migrasi (dari root repo, supaya Makefile menemukan .env)
make migrate-up

# 5. Seed superadmin & data awal
go run ./cmd/seeder

# 6. Jalankan API server (dev: hot-reload dengan air)
make run          # atau: make air
```

## Perintah Berguna

| Perintah | Fungsi |
|---|---|
| `make run` | Jalankan API server (`go run ./cmd/api`) |
| `make air` | Dev server dengan hot-reload |
| `make build` | Build semua package |
| `make fmt` | Format kode |
| `make test` | Jalankan test |
| `make migrate-up` | Jalankan migrasi pending |
| `make migrate-down` | Rollback 1 migrasi |
| `make migrate-create name=xxx` | Buat file migrasi baru |
| `go vet ./...` | Static analysis |

## Struktur

```
backend/
├── cmd/
│   ├── api/          # Entry point server
│   └── seeder/       # Seeder superadmin
├── config/           # Baca env → struct (tidak ada tuning di sini)
├── internal/
│   ├── app/          # Wiring: router, handler, service, repo
│   ├── helper/       # Response, error, audit meta
│   ├── infrastructure/
│   ├── middleware/   # Auth, role, rate limit, captcha
│   ├── module/       # auth, sliders, berita, kegiatan, kontak, pengurus,
│   │                 # role, user, settings, activitylog, dashboard
│   ├── validator/    # Validasi binding → errors[] terstruktur
├── migrations/       # SQL migrasi (goose, immutable — jangan edit yang sudah jalan)
├── public/           # File upload (logo, kegiatan, dsb)
└── seeder/
```

## Variabel Environment

Lihat `.env.example` untuk daftar lengkap + nilai contoh. Catatan penting:

- `JWT_SECRET` — wajib minimal 32 karakter di production (app fail-fast).
- `CAPTCHA_ENABLED=true` mewajibkan `CAPTCHA_SECRET_KEY` terisi (app fail-fast).
- `CAPTCHA_SITE_KEY` — key publik Turnstile, dikirim ke frontend via `GET /settings`.
- `DB_PASS` wajib terisi di production.
- Secret (JWT, DB_PASS, SMTP_PASS, CAPTCHA_SECRET_KEY) JANGAN di-commit.
