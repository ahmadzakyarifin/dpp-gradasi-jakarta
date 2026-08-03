# Tech Stack — dpp-gradasi-jakarta (repo root)

- Repo: `~/Documents/commit-2026/dpp-gradasi-jakarta` — git root, berisi `backend/`, `frontend/`, `bruno-api/`, `docs/`, `docs-final/`, `graphify-out/`.
- Backend: Go 1.25, module `github.com/ahmadzakyarifin/dpp-gradasi/backend`. Gin + GORM (MySQL). Redis-free total (in-memory rate limiter). Gomail untuk email sync.
- 11 modul: auth, user, role, activitylog, dashboard, berita, kegiatan, kontak, pengurus, settings, sliders.
- Pola layer baku (FINAL 2026-08-03): DTO (json+binding) → Entity (tanpa tag) → Model (gorm) → Mapper 2 arah → Handler (tipis) → Service (logic+AUDIT) → Repo (query GORM).
- Frontend: React + Vite + Tailwind + Zustand + React Router (di `frontend/`).
- API contract final: `docs-final/api/*.jsonc`. Postman collection di `bruno-api/`.
- Konfigurasi: `config/*.go` baca env (caarlos0/env). Tanpa `.env` di disk (di-gitignore, ada `.env.example`? — cek). Secret wajib di env.
- Build hijau: `go build ./...` + `go vet ./...` + `gofmt -l .` kosong (verified 2026-08-03).
- Graphify: `graphify update .` dari root repo (AST-only, tanpa LLM karena tanpa GEMINI_API_KEY).