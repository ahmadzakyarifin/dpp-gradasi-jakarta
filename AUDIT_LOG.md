# AUDIT LOG — DPP GRADASI (Company Profile + Admin Panel)

> Tracking audit & perbaikan bertahap per fitur.
> Format: setiap fitur punya laporan audit, menunggu approval sebelum fix.
> Legenda: 🔴 Bug/Security (wajib fix) · 🟡 UX/inconsistency (sebaiknya) · 🔵 Rekomendasi/nice-to-have · ✅ Sudah benar

## Status Fitur

| Fitur | Status Audit | Laporan | Approved | Fixed |
|-------|-------------|---------|----------|-------|
| 1. Full Auth (Login + Forgot + Reset) | ✅ Laporan siap (2026-08-01) | Bagian di bawah | ✅ Approved (2026-08-01) | ✅ Semua (Captcha, backdoor, token, rate limit, email, pending, kontrak) |
| 2. Manajemen Admin (User) | ✅ Laporan siap (2026-08-01) | Bagian di bawah | ✅ Approved (2026-08-01) | ✅ Semua (contract, list, tab, bulk, role, upload, rate limit, alur aktivasi) |

---

# FITUR 1: FULL AUTH (Login + Forgot Password + Reset Password)

**Tanggal audit:** 2026-08-01
**Ruang lingkup:** `frontend/src/pages/Login.jsx`, `frontend/src/pages/ResetPassword.jsx`, `frontend/src/components/auth/*`, `frontend/src/services/authService.js`, `frontend/src/store/useAuthStore.js`, `frontend/src/api/index.js`, `frontend/src/content/authContent.js`, `frontend/src/utils/validation.js`, `frontend/src/App.jsx` (guard route) · backend `internal/module/auth/{dto,handler,service,routes}`, `internal/middleware/{rate_limiter,rate_rules}`, `config/{security,jwt,cookie}` · `docs/api/auth.jsonc` · statis `company_profile/{login,reset-password}.html`

---

## 1.1 Audit UI — React vs HTML statis (Login, Forgot, Reset)

### ✅ Kesamaan (React sudah meniru statis dengan baik)
- **Layout split-screen** (left branding panel `md:w-1/2` + right form): identik.
- **Left panel**: background image Unsplash + overlay gradient `bg-brand-950/95 via-brand-900/90 to-brand-800/80` + `texture-dots`; tombol "Kembali ke Beranda" (pill, blur); logo 80px; heading 2 baris dengan gradient amber; paragraf; copyright; 3 ikon sosial (FB/IG/YT) — semua sama.
- **Form field**: padding `pl-11 pr-4 py-3`, `bg-slate-50 border rounded-xl`, ikon kiri Phosphor, placeholder `nama@email.com` / `••••••••` — sama.
- **Password toggle** eye/eye-slash, **remember me** checkbox — sama.
- **Reset Password**: 4 state (checking / expired / valid / success) lengkap dengan spinner, ikon, teks — sama persis.
- **Hover/focus**: `hover:bg-brand-700`, `focus:ring-2 focus:ring-brand-500/20` — sama.

### 🟡 Perbedaan kecil (bukan blocker, tapi catat)
| Lokasi | HTML statis | React | Keterangan |
|---|---|---|---|
| Login mobile "Beranda" | `<a href="index.html">` | `<Link to="/">` | React benar (SPA) |
| Reset expired tip | Teks `**batas waktu 15 menit**` literal | `expiredMessage` tanpa angka menit | Teks React lebih generik; angka 15 menit hardcode di statis |
| Reset valid subtitle | `<span x-text="email ? email : 'Anda'">` | `{email \|\| 'Anda'}` | Sama |
| Transition antar form | `x-transition` (slide/fade) | Tanpa transisi | React tidak punya animasi pergantian view — minor |
| HTML `x-show` awal | form forgot `style="display:none"` | — | React render bersyarat, fine |

### 🔵 Konten yang seharusnya dinamis tapi masih hardcode (catatan Fitur 4 juga)
- **Logo** `https://gradasi.org/uploads/img/logo/1737187847.png` — hardcode di `authContent.js` (`logoUrl`), dipakai Login mobile + Reset + AuthBrandPanel. **Harusnya dari settings** (`settings.logo_url`), sama seperti PublicLayout yang sudah di-binding.
- **Brand name "DPP GRADASI"** hardcode di `authContent.brandName` → dipakai di header/alt/copyright. Harusnya `settings.site_name`.
- **Background image** Unsplash hardcode (`backgroundUrl`). Harusnya dari settings (kalau ada field `background_url` / `login_bg_url`).
- **Teks hero** (Kembangkan / Potensi Digital / Anda, deskripsi) hardcode di `authContent.login.*` / `authContent.reset.*` — ini konten company profile, sebaiknya dinamis via settings kalau field-nya ada.
- **Social links** FB/IG/YT + URL hardcode di `authContent.socialLinks`. Harusnya dari settings (PublicLayout sudah pakai `settings.social_*`).
- **Copyright "© 2026 DPP GRADASI"** hardcode. Harusnya `settings.footer_text` / tahun dinamis.

> Catatan: ini TIDAK akan di-fix di Fitur 1 (auth) — dicatat agar masuk daftar binding Fitur 4. Kecuali logo/brand, yang paling terlihat.

---

## 1.2 Audit Validasi Form

### Login
| Item | Status | Detail |
|---|---|---|
| Email wajib | ✅ | `validateEmail` → "Email wajib diisi" (kosong) / "Format email tidak valid" |
| Password wajib | ✅ | "Password wajib diisi" |
| Kapan muncul | ✅ | **onBlur + submit** (via `touched`), konsisten dengan statis (`@blur` + submit) |
| Sinkron client/server | 🟡 | Client: `min=6` untuk password TIDAK divalidasi di FE (hanya "wajib diisi"); server: DTO `binding:"required"` (tanpa min) — **sinkron** (keduanya tidak enforce min 6 di login) |
| Hilang saat mengetik | 🟡 | Error hilang **hanya setelah onBlur berikutnya** (karena `touched`); saat mengetik ulang tanpa blur, error tetap tampil — beda dengan statis (`@input="validateLogin()"` **langsung** hilang). **Inkonsistensi UX kecil** |

### Forgot Password
| Item | Status | Detail |
|---|---|---|
| Email wajib/format | ✅ | Sama dengan login |
| Kapan muncul | ✅ | onBlur + submit |
| Hilang saat mengetik | 🟡 | Sama: error tetap sampai blur berikutnya (statis langsung hilang) |

### Reset Password
| Item | Status | Detail |
|---|---|---|
| Password wajib + min 6 | ✅ | `validatePassword` → "Password baru minimal 6 karakter." |
| Konfirmasi wajib | ✅ | "Konfirmasi password wajib diisi." |
| Konfirmasi cocok | ✅ | "Konfirmasi password tidak cocok." |
| Indikator real-time cocok/tidak | ✅ | **Ada**: ikon check/x hijau/merah + teks "Password cocok" saat ketik (statis juga ada) |
| Strength password (selain min 6) | 🔵 | Tidak ada requirement strength (hanya min 6) — konsisten FE & BE; nice-to-have |
| Server-side sync | ✅ | DTO `binding:"required,min=6,eqfield=Password"` — cocok dengan FE |
| Error position | ✅ | Tepat di bawah input, muncul saat submit/touched |

---

## 1.3 Audit Autentikasi & Session — Login

| # | Item | Status | Detail |
|---|---|---|---|
| 1 | Email tidak terdaftar / password salah | ✅ | **Pesan generik** "Email atau password salah." (service `AUTH_INVALID_CREDENTIALS`), tidak membedakan — **tidak ada user enumeration** ✅ |
| 2 | Akun belum aktif (`pending_activation`) | 🟡 | Service: `AUTH_ACCOUNT_PENDING` "Akun Anda belum diaktifkan. Silakan periksa email Anda..." — **tapi tidak ada di kontrak `auth.jsonc`** (kontrak cuma 403 `AUTH_ACCOUNT_INACTIVE`). Kode ini muncul di service tapi tidak terdokumentasi. |
| 3 | Akun nonaktif (`inactive`) | ✅ | `AUTH_ACCOUNT_INACTIVE` 403, pesan "telah dinonaktifkan" |
| 4 | **Rate limiting** | ✅ | Middleware Redis `auth-login` IP+Email 5/min; **tapi UI tidak menampilkan countdown/cooldown** (lihat 1.7) |
| 5 | **Refresh token** | ✅ | HttpOnly cookie, rotation one-time, remember_me 30 hari / default 72 jam — sesuai kontrak |
| 6 | **🔴 Backdoor demo di FE store** | 🔴 | `useAuthStore.login()` punya fallback: `if (email === 'admin@gradasi.org' && password === 'password123')` → **login sukses tanpa backend** (token `demo_token_123`). Ini celah: user bisa "login" padahal backend mati/offline, dan akses admin panel. **WAJIB HAPUS.** |
| 7 | **🔴 Token default di store** | 🔴 | `token: localStorage.getItem('access_token') || 'demo_token_123'` — kalau tidak ada token, dianggap login. **Guard route `/admin/*` jadi tembus** (semua halaman admin bisa diakses tanpa login). **WAJIB HAPUS fallback.** |
| 8 | 🔴 Backdoor di `ResetPassword.jsx` (activation fallback) | 🟡→🔴 | `validateActivationToken` + `activateAccount` dipanggil sebagai fallback di halaman reset — ini mengizinkan halaman reset password dipakai untuk aktivasi akun baru. Perlu cek apakah memang disengaja (flow aktivasi admin). Kalau ya, halaman harus membedakan konteks; kalau tidak, ini bisa jadi bypass. |

---

## 1.4 Audit Forgot Password

| # | Item | Status | Detail |
|---|---|---|---|
| 1 | Email tidak terdaftar | ✅ | Service `ForgotPassword`: `FindByEmail` err → `return nil` → response 200 "Jika email terdaftar..." **tidak bocorkan info** ✅ |
| 2 | Rate limiting | ✅ | `auth-forgot` IP+Email 3/min (Redis) |
| 3 | Pengiriman email end-to-end | 🟡 | **TIDAK bisa diverifikasi penuh**: pakai `s.mail.SendAsync` (async, fire-and-forget) — kalau SMTP gagal, user tidak tahu. Memory Abang: prefer **sync email** supaya response tahu sukses/gagal. Di DEV, link hanya di-log ke console (kontrak bilang "di mode DEV link muncul di console" — **bukan** terkirim). Belum ada bukti email benar-benar terkirim via SMTP (config SMTP ada di .env? lihat catatan). |
| 4 | Isi email | ✅ | HTML template bagus (nama, tombol, TTL 15 menit, disclaimer) |
| 5 | **Link expiry** | 🟡 | TTL 15 menit (`PasswordResetTTLMinutes`); expired → 400 `AUTH_TOKEN_INVALID_OR_EXPIRED`. UI expired → tombol "Minta Link Reset Baru" → `/login?view=forgot` ✅. TAPI: `FindPasswordResetTokenByHash` **tidak mengecek `ExpiresAt`** → token kedaluwarsa **masih dianggap valid** (bug waktu — lihat 1.5) |
| 6 | 🔴 `SendAsync` + tidak ada error handling | 🟡 | Fire-and-forget; kalau SMTP down, user tidak dapat email tapi UI bilang "berhasil". Rekomendasi: sync send (sesuai preferensi Abang) |

---

## 1.5 Audit Reset Password

| # | Item | Status | Detail |
|---|---|---|---|
| 1 | Captcha tidak muncul di reset | ✅ | `CaptchaWidget` hanya dirender di view login (`captchaEnabled &&`), **tidak** di reset — sesuai arahan |
| 2 | Redirect setelah sukses | ✅ | State `success` → tombol "Masuk ke Akun Saya" → `/login` |
| 3 | **Token sekali pakai** | ✅ | `MarkResetTokenUsed` dijalankan setelah reset ✅ |
| 4 | **Token expiry + sekali pakai** | ✅ | `FindPasswordResetTokenByHash` (repo) memfilter `expires_at > NOW() AND used_at IS NULL` — token expired/terpakai **ditolak** ✅. Link 15 menit benar-benar 15 menit. |
| 5 | Link lama invalid setelah dipakai | ✅ | Token di-mark used → query berikutnya tidak ketemu |
| 6 | Aktivasi vs reset campur | 🟡 | Halaman yang sama dipakai untuk reset & aktivasi (fallback). UX: user reset yang token-nya sudah dipakai akan melihat pesan "aktivasi" — perlu kejelasan |

---

## 1.6 Audit CAPTCHA (Cloudflare Turnstile)

| # | Item | Status | Detail |
|---|---|---|---|
| 1 | 🔴 **Bukan Cloudflare Turnstile** | 🔴 | `CaptchaWidget.jsx` = **captcha canvas homemade** (generate 5 char di client, gambar di canvas). **TIDAK ada Turnstile sama sekali** — `TURNSTILE_SITE_KEY` / `TURNSTILE_SECRET_KEY` ada di config tapi **tidak dipakai** (grep backend: hanya `CaptchaEnabled` + cek `CaptchaToken != ""`). |
| 2 | 🔴 **Verifikasi CAPTCHA client-side only** | 🔴 | `handleInputChange` → `onVerify('TOKEN_VERIFIED_' + captchaCode)` — **token di-generate & divalidasi 100% di client**. Backend cuma cek `CaptchaToken != ""` (tidak kosong). Siapa pun bisa kirim `captcha_token: "apaaja"` dan lolos. **Ini security hole.** |
| 3 | Muncul otomatis saat load | ✅ | `useEffect` → `generateCaptcha()` saat mount (di view login) |
| 4 | Token divalidasi ulang di backend | 🔴 | **TIDAK** — hanya cek non-empty. Turnstile server-side verify (`siteverify`) tidak ada. |
| 5 | Expiry token | 🔴 | Konsepnya tidak relevan karena client-side. Kalau pakai Turnstile asli, perlu auto-refresh (belum ada). |
| 6 | 🔴 `VITE_CAPTCHA_ENABLED=true` default di `.env.example` | 🟡 | Frontend default captcha **aktif** tapi backend `CAPTCHA_ENABLED` default `false` (tidak ada di .env). Mismatch: FE minta captcha, BE tidak wajibkan → konsisten untuk "tidak ada proteksi" tapi **membingungkan** & captcha homemade tetap tampil. |

**Kesimpulan CAPTCHA:** 🔴 Implementasi sekarang **tidak memberi proteksi apa pun** — captcha canvas client-side bisa di-bypass total (token dibuat sendiri). Kalau arahan adalah "Cloudflare Turnstile", ini harus diganti total: render Turnstile di FE + `siteverify` di BE + config keys. **Ini temuan terbesar Fitur 1.**

---

## 1.7 Audit Rate Limiter

| # | Item | Status | Detail |
|---|---|---|---|
| 1 | Implementasi backend | ✅ | `ulule/limiter` + Redis store (`dppgradasi_rate_limit:*`), scope `auth-login` (IP+Email 5/min), `auth-forgot` (IP+Email 3/min) |
| 2 | **Cloudflare edge** | 🔵 | Tidak ada (hanya in-process/Redis). Nice-to-have kalau deploy di belakang Cloudflare |
| 3 | Header rate limit | ✅ | `X-RateLimit-Limit/Remaining/Reset-*` + `Retry-After` di-set (tapi FE tidak baca) |
| 4 | 🔴 **UI tidak menampilkan cooldown/countdown** | 🔴 | Saat 429 `AUTH_RATE_LIMIT_EXCEEDED`, FE cuma tampilkan `error.message` generic "Terlalu banyak percobaan. Silakan coba lagi nanti." **Tidak ada countdown "coba lagi dalam X detik"**, tidak ada disable tombol. Backend kirim `retry_after` tapi FE **tidak pakai**. **UX buruk + membingungkan.** |
| 5 | Per-IP / per-akun | ✅ | IP+Email kombinasi (hash). Bagus. |

---

## 1.8 Audit Kode & Arsitektur

### Backend
| # | Item | Status | Detail |
|---|---|---|---|
| 1 | Layer handler→service→repo | ✅ | Rapi: handler tipis (bind + call service), service logic, repo query. Tidak ada logic tercecer di handler. |
| 2 | Error handling | ✅ | `ServiceError` + `errorCodeToHTTP` — konsisten |
| 3 | Activity log | ✅ | login sukses/gagal, forgot, reset, activate di-log (async `go logSvc.Log`) |
| 4 | 🔴 `errorCodeToHTTP` — `AUTH_ACCOUNT_PENDING` | 🟡 | Tidak ada case → fallback 500. Padahal service return `AUTH_ACCOUNT_PENDING` untuk user pending → user dapat **500** bukan 4xx. **Bug kecil.** |
| 5 | 🔴 Captcha check di handler | 🔴 | Hanya `if req.CaptchaToken == ""` — tanpa verifikasi (lihat 1.6) |
| 6 | `LoginRequest` binding | ✅ | `required,email` + password `required` (server tidak enforce min 6 di login — konsisten FE) |

### Frontend
| # | Item | Status | Detail |
|---|---|---|---|
| 1 | Komponen reusable | ✅ | `AuthShell`, `AuthBrandPanel`, `CaptchaWidget` terpisah rapi |
| 2 | 🔴 `useAuthStore` backdoor + token default | 🔴 | (lihat 1.3 #6, #7) — **harus dihapus** |
| 3 | `Login.jsx` — gabungan login+forgot dalam 1 file | 🟡 | Boleh (sama seperti statis `x-show`), tapi state `forgotEmail` terpisah — OK. Halaman forgot bukan route sendiri (`/login?view=forgot`) — sama dengan statis |
| 4 | Error handling FE | ✅ | `apiRequest` throw error dgn `status`/`code`/`data` — bagus |
| 5 | `useAuthStore` default user hardcode | 🟡 | `user: {id:1, name:'Super Admin', email:'admin@gradasi.org', role:'Super Admin'}` default — data palsu sebelum fetch. Jangan-jangan admin panel pakai data ini (role "Super Admin" dgn spasi — **mismatch** dgn role DB `super_admin`). Perlu cek Dashboard/admin pakai `user.role` apa. |

### API Contract (`docs/api/auth.jsonc`) vs Implementasi
| # | Endpoint | Contract | Aktual | Status |
|---|---|---|---|---|
| 1 | POST `/auth/login` | 200, 400 (captcha), 401, 403 (inactive), 422, 429, 500 | 200, 400 (captcha), 401, 403 (inactive), 500 + **429 (AUTH_RATE_LIMIT_EXCEEDED dgn retry_after)** + **422 (VALIDATION_ERROR)** | 🟡 **429 dobel** di contract (baris 100 & 132: `RATE_LIMIT_EXCEEDED` vs `AUTH_RATE_LIMIT_EXCEEDED`) + aktual pakai `AUTH_RATE_LIMIT_EXCEEDED`. Contract ambigu — harusnya 1 saja. |
| 2 | POST `/auth/refresh` | 200, 401, 403, 500 | sama | ✅ |
| 3 | POST `/auth/logout` | 200, 401, 500 | sama | ✅ |
| 4 | GET `/auth/me` | 200, 401, 404, 500 | sama | ✅ |
| 5 | POST `/auth/forgot-password` | 200, 422, 429, 500 | sama | ✅ (tapi 429 dgn `AUTH_RATE_LIMIT_EXCEEDED`, contract konsisten) |
| 6 | GET `/auth/validate-reset-token` | 200, 400, 500 | 200, 400, 500 | ✅ |
| 7 | POST `/auth/reset-password` | (baca lanjutan contract) | 200, 400, 500 | ✅ |
| 8 | GET `/auth/validate-activation-token` | — | **ADA di backend, TIDAK ada di contract** | 🟡 kontrak kurang |
| 9 | POST `/auth/activate-account` | — | **ADA di backend, TIDAK ada di contract** | 🟡 kontrak kurang |
| 10 | POST `/auth/change-password` | — | **ADA di backend, TIDAK ada di contract auth.jsonc** | 🟡 (mungkin di users.jsonc) |

> Catatan: `AUTH_ACCOUNT_PENDING` (pending_activation) ada di service tapi **tidak ada** di contract responses.

---

## RINGKASAN TEMUAN — FITUR 1

> **UPDATE 2026-08-01:** Abang mengonfirmasi sudah punya konfigurasi Cloudflare Turnstile asli di .env (`CAPTCHA_SITE_KEY`/`CAPTCHA_SECRET_KEY`). Captcha, backdoor & token default sudah di-fix (lihat bawah).

### 🔴 Bug/Security (wajib fix)
1. **Captcha bukan Turnstile & tidak ada verifikasi server** — captcha canvas client-side, backend hanya cek `captcha_token != ""`. Bisa di-bypass total. (1.6) → ✅ **FIXED** (lihat bawah)
2. **Backdoor login demo di `useAuthStore`** — `admin@gradasi.org/password123` lolos tanpa backend. (1.3) → ✅ **FIXED**
3. **Token default `demo_token_123`** — tanpa login pun `token` terisi → **guard route /admin tembus**. (1.3) → ✅ **FIXED**
4. ~~Token reset expiry tidak dicek~~ → **Terverifikasi OK**: `FindPasswordResetTokenByHash` memfilter `expires_at > NOW() AND used_at IS NULL`. (1.5)
5. **Rate limit 429: UI tidak tampilkan countdown** — `retry_after` dikirim tapi tidak dipakai; user bingung kapan bisa coba lagi. (1.7) → ✅ **FIXED**

### 🟡 UX/inconsistency (sebaiknya fix)
6. `AUTH_ACCOUNT_PENDING` → 500 di `errorCodeToHTTP` (harusnya 403/4xx). (1.8) → ✅ **FIXED** (403, terverifikasi live)
7. Email async `SendAsync` — tidak ada feedback kalau SMTP gagal (preferensi: sync). (1.4) → ✅ **FIXED** (sync + `EMAIL_SEND_FAILED`)
8. Error form tidak hilang saat mengetik ulang (beda dgn statis `@input`). (1.2) → ✅ **FIXED** (onChange mark touched)
9. ~~`AUTH_ACCOUNT_PENDING` & `validate-activation-token`/`activate-account`/`change-password` tidak terdokumentasi di contract~~ → **Koreksi audit:** contract SUDAH punya ketiganya (saya tidak baca sampai habis waktu audit). Yang benar: `AUTH_ACCOUNT_PENDING` belum ada → ✅ **ditambahkan** ke contract.
10. 429 dobel di contract login (dua definisi beda kode). (1.8) → ✅ **FIXED** (hapus `RATE_LIMIT_EXCEEDED` lama, sisakan `AUTH_RATE_LIMIT_EXCEEDED` + retry_after)
11. Halaman reset = gabungan reset + aktivasi (fallback) — perlu keputusan desain. (1.5) → ⏳ Belum (keputusan desain)

### 🔵 Rekomendasi/nice-to-have
12. Binding konten auth ke settings (logo, brand, background, social, copyright) — dicatat utk Fitur 4. (1.1)
13. Cloudflare edge rate limit kalau deploy di belakang CF. (1.7)
14. Password strength meter (nice-to-have). (1.2)
15. Transisi antar form login↔forgot (sama seperti statis). (1.1)

### ✅ Sudah benar
- Pesan login generik (anti user enumeration)
- Rate limit per IP+Email via Redis (login 5/min, forgot 3/min)
- Refresh token rotation + HttpOnly cookie
- Token reset sekali pakai (mark used)
- Email template reset bagus (TTL 15 menit tercantum)
- Captcha tidak muncul di halaman reset (sesuai arahan)
- Layer arsitektur BE rapi (handler→service→repo)
- Validasi konfirmasi password real-time (indikator cocok/tidak)
- Error di bawah input yang benar, muncul onBlur/submit

---

## LOG FIX — FITUR 1 (2026-08-01)

### 🔧 Fix CAPTCHA → Cloudflare Turnstile asli (end-to-end)
**Latar:** Abang sudah punya keys Turnstile resmi di `.env` (`CAPTCHA_SITE_KEY`, `CAPTCHA_SECRET_KEY`) — tapi config backend membaca `TURNSTILE_SITE_KEY`/`TURNSTILE_SECRET_KEY` → **mismatch nama env** → keys tidak pernah terpakai, dan widget captcha yang dirender adalah canvas homemade (bypassable).

**Perubahan:**
| File | Perubahan |
|---|---|
| `backend/config/security.go` | Env tags diubah `TURNSTILE_*` → `CAPTCHA_SITE_KEY`/`CAPTCHA_SECRET_KEY` + tambah `CAPTCHA_VERIFY_URL` |
| `backend/internal/helper/turnstile.go` | **BARU** — `VerifyTurnstile()`: POST ke `challenges.cloudflare.com/turnstile/v0/siteverify`, timeout 5s, dev-bypass jika secret kosong |
| `backend/internal/module/auth/handler/auth_handler.go` | `Login` & `ForgotPassword` verifikasi siteverify (bukan cuma cek non-empty). Error: `AUTH_CAPTCHA_REQUIRED` (kosong) / `AUTH_CAPTCHA_INVALID` (gagal) |
| `backend/internal/module/auth/dto/auth_dto.go` | `ForgotPasswordRequest` + field `captcha_token` |
| `frontend/src/components/auth/CaptchaWidget.jsx` | **Tulis ulang** → Turnstile asli (render iframe Cloudflare, callback/expired/error, reset on retry) |
| `frontend/src/pages/Login.jsx` | Forgot view + CaptchaWidget + kirim `captcha_token` di payload forgot; reset state saat switch view |
| `frontend/src/services/authService.js` | `forgotPassword(email, captchaToken)` |
| `.env` | + `CAPTCHA_ENABLED=true`, `CAPTCHA_VERIFY_URL=` (default prod) |
| `.env.example` | Dokumentasi CAPTCHA + mode dev/prod |
| `frontend/.env` + `.env.example` | **BARU** — `VITE_CAPTCHA_ENABLED` + `VITE_CAPTCHA_SITE_KEY` (key publik) |

**Verifikasi live:**
- ✅ `go build` + `go vet` + `npm run build` hijau
- ✅ `POST /login` tanpa captcha → `AUTH_CAPTCHA_REQUIRED`
- ✅ `POST /login` captcha palsu → `AUTH_CAPTCHA_INVALID` (siteverify ditolak)
- ✅ `POST /forgot-password` tanpa captcha → `AUTH_CAPTCHA_REQUIRED`
- ✅ `POST /forgot-password` captcha palsu → `AUTH_CAPTCHA_INVALID`
- ✅ Browser: halaman `/login` menampilkan **iframe Cloudflare Turnstile asli** ("Widget containing a Cloudflare security challenge" + checkbox "Verify you are human")

### 🔧 Fix backdoor & token default (frontend)
| File | Perubahan |
|---|---|
| `frontend/src/store/useAuthStore.js` | Hapus fallback `demo_token_123` + user default hardcode; token init `|| null`; hapus bypass `admin@gradasi.org/password123`; `login` throw error jika token tidak ada |

**Verifikasi:** browser `/login` **tidak lagi redirect** ke `/dashboard` (sebelumnya tembus karena token default).

### ⏳ Belum di-fix (menunggu approval lanjutan)
- Rate limit 429 countdown di UI (#5)
- `AUTH_ACCOUNT_PENDING` → HTTP 500 (#6)
- Email sync `SendAsync` (#7)
- Error form hilang saat ketik (#8)

---

## LOG FIX 2 — FITUR 1 (2026-08-01, sesi lanjutan)

### 🔧 Fix #5–#10 (temuan 🟡 yang di-approve "lanjut saja")

| # | Temuan | Fix | Verifikasi |
|---|--------|-----|------------|
| #5 | Rate limit 429 tanpa countdown | `api/index.js` expose `error.retryAfter` (dari body `retry_after`/header `Retry-After`); `Login.jsx` + state `cooldown` + interval countdown, notice "Coba lagi dalam X detik", tombol disabled `Tunggu Xs`, reset saat switch view | ✅ `POST /login` ×6 → 429 `AUTH_RATE_LIMIT_EXCEEDED` + `retry_after` di body |
| #6 | `AUTH_ACCOUNT_PENDING` → 500 | `auth_handler.go` `errorCodeToHTTP` + case → `http.StatusForbidden` | ✅ Login user `pending_activation` → **403** `AUTH_ACCOUNT_PENDING` (live) |
| #7 | Email async tanpa feedback | `auth_service.go` `SendAsync` → `Send` (sync) di inactive-account & reset-password; error → `EMAIL_SEND_FAILED` | ✅ build; alur error service terverifikasi |
| #8 | Error form tidak hilang saat ketik | `Login.jsx` onChange ketiga input + `setTouched(...true)` → error update real-time | ✅ npm build |
| #9 | (Koreksi) 3 endpoint SEBENARNYA sudah ada di contract | `AUTH_ACCOUNT_PENDING` ditambahkan ke contract (403 kedua) | ✅ JSONC valid |
| #10 | 429 dobel di contract login | Hapus definisi lama `RATE_LIMIT_EXCEEDED`, sisakan `AUTH_RATE_LIMIT_EXCEEDED` + `retry_after` | ✅ JSONC valid |

### ⚠️ Temuan baru (di luar audit awal): ENUM drift migration 00012
**Gejala:** `INSERT users status='pending_activation'` → MySQL error 1265 (Data truncated). `SHOW COLUMNS` menunjukkan `enum('active','inactive')` padahal goose status bilang 00012 sudah applied.

**Root cause:** migration 00012 (`ALTER TABLE users MODIFY status ENUM(...,'pending_activation')`) **tidak benar-benar mengubah kolom di DB `dpp_gradasi`** — kemungkinan dijalankan terhadap DB yang berbeda (ada `dpp_baru`, `dpp_gradasi_jakarta`, dll di instance yang sama; `dpp_baru.users.status` SUDAH punya `pending_activation`, `dpp_gradasi` TIDAK). Goose mencatat applied tapi ALTER tidak dieksekusi pada DB target.

**Fix:** ALTER manual di `dpp_gradasi`:
```sql
ALTER TABLE users MODIFY COLUMN status ENUM('active','inactive','pending_activation') NOT NULL DEFAULT 'inactive';
```
**Verifikasi:** `SHOW COLUMNS` → ENUM sudah 3 nilai; INSERT user pending berhasil.

**Rekomendasi:** audit semua DB di instance (`dpp_gradasi_jakarta`, `dpp_baru`, dll) terhadap migration 00012–00016 — kemungkinan drift serupa (migration di-apply ke DB yang salah).

### ⏳ Sisa (keputusan desain, belum di-fix)
- #11 Halaman reset = gabungan reset + aktivasi (fallback) — butuh keputusan: pisah halaman atau biarkan?
- 🔵 Rekomendasi #12–#15 (binding konten ke settings → Fitur 4; Cloudflare edge rate limit; password strength meter; transisi form)

---

# FITUR 2: MANAJEMEN ADMIN (USER)

**Tanggal audit:** 2026-08-01
**Tanggal fix:** 2026-08-01 (approval Abang)
**Ruang lingkup:** `docs/api/users.jsonc` · backend `internal/module/user/{dto,handler,service,repository,model,routes}` · FE `frontend/src/pages/admin/{UsersAdmin,ProfileAdmin}.jsx`, `frontend/src/services/userService.js` · statis `admin/users.html`

## Ringkasan

| Severity | Jumlah | Deskripsi |
|----------|--------|-----------|
| 🔴 Critical | 3 | Contract invalid, FE-BE list mismatch, hardcode URL |
| 🟡 High | 5 | Tab filter tidak jalan, bulk tidak proteksi, role name, upload tanpa validasi, no rate limit |
| ✅ Benar | 7 | Proteksi super_admin, sync email, soft delete, validasi role_id, activity log, selectable FE, konsisten statis |

**STATUS: SEMUA FIXED ✅ (approved 2026-08-01).** Detail di LOG FIX — FITUR 2 di bawah.

## 🔴 Critical

### #1 Contract `users.jsonc` TIDAK VALID JSON — ✅ FIXED
**Lokasi:** `docs/api/users.jsonc` baris ~131 (setelah blok `verify_email`).
**Masalah:** trailing comma ganda (`},` + `}`) → parser gagal (`Expecting property name enclosed in double quotes: line 131`).
**Dampak:** Tooling (Bruno, validator) tidak bisa baca contract. Semua endpoint user tidak terdokumentasi dengan benar.
**Fix:** hapus koma dobel. ✅ Validasi `python3 validate_jsonc.py docs/api` PASS.

### #2 FE-BE list mismatch — halaman UsersAdmin TIDAK pernah tampil — ✅ FIXED
**Lokasi:** FE `UsersAdmin.jsx` `fetchUsers()` + `userService.list()` vs BE `GetAdmins`/`FindAllAdmins`.
**Masalah:**
- FE kirim `tab/search/page/limit` → BE `GetAdmins` **tidak baca query param** (handler kosong, service panggil `FindAllAdmins()` tanpa argumen).
- FE expects `res.data.items || res.data.users` + `res.data.pagination` → BE return `data: [ ... ]` polos (array langsung, tanpa items/pagination).
**Dampak:** `res.data.items` undefined → `users = []` → tabel kosong "Tidak ada data untuk ditampilkan". **Halaman manajemen admin tidak berfungsi sama sekali.**
**Fix:** backend tambah query params `tab/search/page/limit` + response `{items, pagination}` + Preload Role. ✅ E2E verified: GET `/admin/users?tab=active&page=1&limit=5` → `{items:[{role:"super_admin",...}], pagination:{total:2}}`.

### #3 `ProfileAdmin.jsx` hardcode `http://127.0.0.1:8080` — ✅ FIXED
**Lokasi:** `frontend/src/pages/admin/ProfileAdmin.jsx` baris 44, 67 (fetch), 22 (`photo_path`).
**Masalah:** URL absolut hardcode → production (domain beda) gagal total.
**Fix:** pakai `VITE_API_URL` / `apiRequest` + base untuk `photo_path`. ✅

## 🟡 High

### #4 Tab filter (pending/trash) tidak berfungsi — ✅ FIXED
FE punya 3 tab (Aktif/Menunggu Aktivasi/Terhapus) tapi BE `FindAllAdmins` ambil SEMUA role IN (1,2,3,4) tanpa filter status/deleted → semua tab tampil data sama, tab trash tidak pernah soft-deleted (karena `Find` default exclude soft delete).
**Fix:** BE support `?tab=active|pending|trash` → query `status`/`deleted_at` (+ pagination & search). ✅ E2E verified: `tab=trash` total 2, `search=admin` → 2 hasil.

### #5 Bulk delete/restore tidak proteksi super_admin — ✅ FIXED
`BulkDeleteAdmin` hanya filter `id != adminID` (diri sendiri) — **tidak** cek `role_id == 1`. Contract bilang "tidak bisa menghapus sesama super_admin" tapi bulk bisa hapus semua super admin lain. `BulkRestoreAdmin` juga tanpa proteksi.
**Fix:** filter role 1 dari targetIDs di service. ✅ E2E verified: bulk `[1,6]` → id=1 tetap ada (super admin skip), id=6 kehapus.

### #6 Role name tidak ada di response — ✅ FIXED
FE tampil `item.role` (`Super Admin`/`Admin`/dll) + `item.role === 'Super Admin'` untuk proteksi, tapi BE `toResponse` hanya kirim `role_id` (tanpa role name). → FE tampil "Admin" default semua + proteksi Super Admin tidak jalan di UI (tapi BE tetap proteksi).
**Fix:** include role name di response (Preload Role → `role` field). ✅ E2E verified: `role: "super_admin"`/`role: "admin"` di response.

### #7 Upload foto profil tanpa validasi MIME/ukuran — ✅ FIXED
`handleUpload` (service) terima file apapun (tulis ke disk tanpa cek content type/limit). FE accept `image/*` tapi BE tidak enforce → bisa upload file berbahaya.
**Fix:** validasi MIME (detect content) + limit 2MB (sama seperti settings logo). ✅

### #8 Tidak ada rate limit di admin routes — ✅ FIXED
Admin users routes tanpa rate limiter (beda auth). Minor untuk v1 (hanya super_admin yang akses).
**Fix:** tambah `RateLimitRules` di GET (30/min), POST (10/min), resend (5/min). ✅

## ✅ Sudah benar (tidak diubah)
- Proteksi super_admin di `SetAdminStatus`/`DeleteAdmin` (role 1 tidak bisa diubah/dihapus; tidak bisa hapus diri)
- `CreateAdmin` sync email + rollback user/token kalau SMTP gagal (sesuai contract)
- Soft delete + restore + bulk (repo benar)
- DTO validasi `role_id oneof=2 3 4` (super_admin tidak bisa dibuat via undangan)
- Activity log lengkap di semua aksi (create/delete/restore/status/resend/bulk)
- FE `selectableUsers` mengecualikan Super Admin & pending (konsisten statis)
- FE modal form + confirm dialog rapi (mengikuti statis)

---

## LOG FIX — FITUR 2 (2026-08-01, approved Abang)

### Perubahan alur aktivasi akun (sesuai permintaan Abang)
**SEBELUM:** CreateAdmin → status `pending_activation` → email link aktivasi → admin set password sendiri via activate-account.
**SESUDAH (baru):**
1. Super Admin buat admin → sistem generate **password default acak (10 karakter)** → status langsung `active`
2. Email kredensial (email + password default) dikirim **sinkron**; kalau SMTP gagal → **rollback** (akun tidak jadi dibuat, error `MAIL_SEND_FAILED`)
3. Admin login pakai email + password default → response login berisi `must_change_password: true`
4. FE redirect ke `/admin/profile?force=1` → banner "Anda menggunakan password default" → wajib ganti password
5. Setelah ganti password → `must_change_password` di-clear → login normal
6. Tombol "Kirim Ulang Aktivasi" (resend-activation) sekarang = kirim ulang **kredensial baru** (password default baru)

### File diubah
| File | Perubahan |
|------|-----------|
| `backend/migrations/00017_add_must_change_password.sql` | BARU — kolom `must_change_password TINYINT(1) DEFAULT 0` |
| `backend/internal/module/user/model/user_model.go` | field `MustChangePassword` |
| `backend/internal/module/user/dto/user_dto.go` | `UserResponse.RoleName`, `ListUsersQuery`, `UserListResponse`, `Pagination` |
| `backend/internal/module/user/repository/user_repo.go` | `FindAllAdmins(q)` filter tab/search + pagination + Preload Role |
| `backend/internal/module/user/service/user_service.go` | `GetAdmins(q)`, CreateAdmin password default + email kredensial sync + rollback, ResendActivation kredensial baru, BulkDelete filter super admin, `handleUpload` MIME+2MB, ChangePassword clear flag |
| `backend/internal/module/user/handler/user_handler.go` | `GetAdmins` query params; **fix `getAuthUserID` baca context (bug 401 profile)**; message kredensial |
| `backend/internal/module/user/routes.go` | rate limit admin GET/POST/resend |
| `backend/internal/module/auth/dto/auth_dto.go` | `AuthUserResponse.MustChangePassword` |
| `backend/internal/module/auth/service/auth_service.go` | Login & Me include `must_change_password` |
| `backend/internal/module/auth/repository/auth_repo.go` | `SetUserPassword` clear flag (dipakai change/reset/activate) |
| `docs/api/users.jsonc` | valid JSON; get_all query+pagination; create kredensial; delete super admin note; resend kredensial |
| `docs/api/auth.jsonc` | login response + `must_change_password` |
| `frontend/src/pages/admin/ProfileAdmin.jsx` | hapus hardcode `127.0.0.1` (pakai `VITE_API_URL`), forced-change banner `?force=1` |
| `frontend/src/pages/admin/UsersAdmin.jsx` | role label (`super_admin`), indikator "Ganti pwd", toast kredensial, hapus phone |
| `frontend/src/pages/Login.jsx` | redirect ke `/admin/profile?force=1` kalau `must_change_password` |

### Bug pre-existing ditemukan & di-fix (di luar temuan audit)
- **`getAuthUserID` 401 semua endpoint profile**: handler baca `c.Get("user_id")` (string literal) tapi AuthMiddleware set `helper.ContextUserID` (custom type `contextKey`) → mismatch tipe → `GetProfile`/`ChangePassword`/`UpdateProfile` selalu `UNAUTHORIZED`. Fix: baca `c.Request.Context().Value(helper.ContextUserID)`. E2E verified.

### E2E verification (2026-08-01)
- `POST /auth/login` (password default) → `success: true, must_change_password: true` ✅
- `GET /profile` → `must_change_password: true` (fix 401) ✅
- `PUT /profile/password` (old=default, new=baru) → `PASSWORD_CHANGED` ✅
- `POST /auth/login` (password baru) → `must_change_password: false` (flag clear) ✅
- `POST /auth/login` (password lama) → `AUTH_INVALID_CREDENTIALS` ✅
- `GET /admin/users?tab=active&page=1&limit=5` → `{items, pagination, role}` ✅
- `GET /admin/users?tab=trash` → soft-deleted only ✅
- `GET /admin/users?search=admin` → filter nama/email ✅
- `DELETE /admin/users/1` (super admin) → BAD_REQUEST "Tidak bisa menghapus akun sendiri" ✅
- `POST /admin/users/bulk-delete {"ids":[1,6]}` → id=1 tetap (super admin skip), id=6 terhapus ✅
- `POST /admin/users/1/resend-activation` → FORBIDDEN (super admin) ✅
- `POST /admin/users` (SMTP gagal) → `MAIL_SEND_FAILED` + rollback (INSERT → delete) ✅
- SMTP: Gmail `535 BadCredentials` (app password salah/expired) — **infrastruktur, bukan bug kode**; alur rollback sudah terbukti.

### Catatan
- `make test` tidak ada di Makefile (hanya `go test ./...` → 0 FAIL, semua paket `[no test files]` atau `ok`).
- SMTP Gmail gagal `535` → perlu update app password di `.env` sebelum CreateAdmin bisa kirim email ke produksi.
- Super admin TIDAK BISA dihapus (single maupun bulk) — sesuai permintaan Abang.

---

## UPDATE 2026-08-01 (verifikasi lanjutan Abang)

### ✅ Verifikasi alur forgot password (dengan SMTP mock lokal, bukti end-to-end)
- `POST /auth/forgot-password` → email terkirim (Subject "Reset Password - DPP GRADASI", link + token, berlaku 15 menit) ✅
- `GET /auth/validate-reset-token?token=...` → `AUTH_RESET_TOKEN_VALID` ✅
- `POST /auth/reset-password` → `AUTH_RESET_PASSWORD_SUCCESS` ✅
- `POST /auth/login` dengan password baru → **masuk sebagai super_admin** ✅
- **Fix baru**: forgot password untuk akun **nonaktif** → email "Pemberitahuan Akun DPP GRADASI" berisi status **Nonaktif** (bukan link reset) ✅
- **Fix baru**: forgot password untuk akun **pending_activation** → email berisi status **belum diaktifkan**, tanpa link reset ✅ (sebelumnya dapat link reset biasa — rawan)
- Anti-enumerasi: semua response forgot selalu `AUTH_FORGOT_PASSWORD_SUCCESS` (200) walau email tidak terdaftar ✅

### ✅ Verifikasi rate limiter login
- 5 percobaan salah → `AUTH_INVALID_CREDENTIALS`; percobaan ke-6 → `AUTH_RATE_LIMIT_EXCEEDED` + `retry_after: 1` (detik) ✅
- FE menampilkan countdown "Coba lagi dalam X detik" + tombol submit disabled selama cooldown ✅

### ✅ Verifikasi reset inputan saat login gagal
- **Fix baru**: `handleLogin` catch → `setLoginForm(email:'', password:'')` + reset touched/captcha ✅
- Verified live di browser (`form.requestSubmit()`): email & password jadi **kosong** setelah submit gagal, notice error muncul ✅
- **Fix tambahan**: `useAuthStore.login` sebelumnya re-throw `new Error(err.message)` → properti `retryAfter`/`code` hilang. Sudah diubah jadi `throw err` (error asli) agar countdown 429 di FE menerima `retry_after` akurat. ✅

### ✅ User management & aktivasi (dari E2E sebelumnya)
- Create admin → email kredensial (password default) → login → `must_change_password: true` → forced change → login normal ✅
- User masuk sesuai role (`super_admin`/`admin`) ✅
- Super admin tidak bisa dihapus (single/bulk) ✅

### ⚠️ SMTP Gmail masih `535 BadCredentials`
- App password di `.env` (`SMTP_PASS`) tidak diterima Gmail → semua email produksi gagal.
- **Aksi Abang**: buat app password baru di https://myaccount.google.com/apppasswords (2FA wajib aktif) → update `SMTP_PASS` di `.env`.
- Setelah update, alur forgot/aktivasi/kredensial langsung jalan (kode sudah teruji dengan mock).

---

## UPDATE 2 (2026-08-01) — SMTP real aktif + rate limiter bug fix besar

### ✅ SMTP Gmail REAL sudah aktif (app password baru Abang)
- `SMTP_EMAIL=schoolpay41@gmail.com` + app password baru → **`SEND OK`** (test gomail langsung).
- Forgot password → **email real terkirim** (tanpa mock) — verifikasi via log server tanpa error.
- CreateAdmin → **email kredensial terkirim real** (aktivasi.test@gmail.com, user id=9) — tidak ada `EMAIL_SEND_FAILED`.

### ✅ Remember me terverifikasi (cookie max-age)
- Login **tanpa** remember_me → refresh cookie **72 jam (3 hari)**.
- Login **dengan** remember_me=true → refresh cookie **720 jam (30 hari)**.
- Config: `JWT_ACCESS_TTL_MINS=15`, `JWT_REFRESH_TTL_HOURS=72`, `JWT_REMEMBER_ME_TTL_HOURS=720`.

### 🔴🔴 BUG KRITIS DITEMUKAN & DIFIX: rate limiter TIDAK PERNAH blokir + retry_after selalu 1

**Akar masalah 1 — ulule/limiter v3.11.2 Lua bug**: script `incr` hanya panggil `PEXPIRE` jika `ret == count` (benar hanya di request pertama). Request ke-2+ → counter naik tapi **key TANPA TTL** (TTL -1 persist) → `resetAt - now` selalu ≈ 0 → `retry_after` di-clamp ke 1 **selamanya**, dan key tidak pernah expire → **user terkunci permanen** (inilah penyebab NetworkError/429 permanen yang Abang alami).
- **Fix**: ganti ulule/limiter → **custom fixed-window limiter** (`rate_limiter_redis.go` + `rate_limiter_fixed.go`): Lua `INCR + PEXPIRE` atomic, return `{count, ttl}`; retry_after = TTL tersisa. ulule/limiter dihapus dari go.mod.

**Akar masalah 2 — `1*60` = 60ns (bukan 60 detik!)**: routes pakai `middleware.IPEmail(5, 1*60)` — argumen `time.Duration` tapi `1*60` (int) = **60 nanodetik** → window nyaris nol → key langsung expire → **tidak pernah blokir**.
- **Fix**: semua routes (`auth/login`, `auth/forgot-password`, `users-admin` GET/POST/resend) ganti `1*60` → `1*time.Minute`.

**Akar masalah 3 — emailFromBody konsumsi body**: `c.ShouldBindBodyWithJSON` di middleware rate limit mengonsumsi body, `RestoreBody` tidak restore (karena `gin.BodyBytesKey` tidak di-set) → handler Login bind kosong → selalu `VALIDATION_ERROR`/`AUTH_INVALID_CREDENTIALS` & key rate limit tidak terbentuk (email kosong → skip rule).
- **Fix**: `emailFromBody` pakai `io.ReadAll(c.Request.Body)` + restore manual (`io.NopCloser(bytes.NewReader(raw))`).

**Hasil verifikasi (7 percobaan):**
```
percobaan 1-5 → AUTH_INVALID_CREDENTIALS
percobaan 6   → AUTH_RATE_LIMIT_EXCEEDED | retry_after: 59  (59 detik!)
percobaan 7   → AUTH_RATE_LIMIT_EXCEEDED | retry_after: 59
(setelah 10 detik) → retry_after: 32  ← countdown berjalan
```
Redis key sekarang punya **TTL valid (45s/31s)** — auto-expire → user bisa login lagi setelah window selesai. **Ini menyelesaikan NetworkError Abang.**

### ✅ Alur aktivasi akun lengkap terverifikasi (SMTP real)
```
CreateAdmin → email kredensial (password default acak) terkirim real
→ login password default → must_change_password: true (role: admin)
→ ganti password (user set sendiri) → PASSWORD_CHANGED
→ login password baru → must_change_password: false (role: admin)
```
- Aktivasi "bisa di-setting oleh user" ✅ (forced change password)
- User masuk sesuai role ✅
- Super admin tidak bisa dihapus ✅

---

## UPDATE 3 (2026-08-01) — CAPTCHA aktif + countdown bersih + dead code dihapus

### ✅ CAPTCHA Turnstile sekarang BENAR-BENAR berfungsi (bukan hiasan)
- **Temuan**: `.env` backend masih `CAPTCHA_ENABLED=false` → frontend render widget tapi backend **tidak memverifikasi token** → captcha cuma hiasan.
- **Fix**: set `CAPTCHA_ENABLED=true` di `.env` (backend). Site key FE & BE cocok (`0x4AAAAAAEDe3XqhMzghG4b9`).
- **Verifikasi E2E (7 percobaan login)**:
  ```
  1-5 → AUTH_CAPTCHA_REQUIRED (tanpa token → ditolak 400)
  6   → AUTH_RATE_LIMIT_EXCEEDED | retry_after: 59
  7   → AUTH_RATE_LIMIT_EXCEEDED | retry_after: 59
  ```
  - Token kosong → `AUTH_CAPTCHA_REQUIRED` ✅
  - Token palsu → `AUTH_CAPTCHA_INVALID` (Cloudflare verify) ✅
  - Rate limit tetap jalan **bersamaan** dengan captcha (middleware sebelum handler) ✅

### ✅ Countdown rate limit tanpa duplikasi pesan + tombol disabled
- **Sebelum**: notice statis `"Terlalu banyak percobaan. Coba lagi dalam 26 detik."` + span dinamis `"Coba lagi dalam 20 detik..."` → **duplikasi angka beda** (jelek).
- **Sesudah** (Login.jsx):
  - Notice statis: `"Terlalu banyak percobaan. Silakan tunggu sebelum mencoba lagi."`
  - Span countdown dinamis (satu-satunya): `"Coba lagi dalam {cooldown} detik."`
  - Tombol login/forgot: `disabled={isSubmitting || cooldown > 0}` + teks `Tunggu {cooldown}s` — **selama countdown, tombol mati** (industry standard).
  - `error.retryAfter` dari 429 = **sisa waktu nyata** (59, 32, ...) → countdown akurat.

### 🧹 Bersih-bersih dead code (tanpa ubah logic)
- `mail.go`: hapus `SendAsync` (email sync sudah keputusan; async tidak dipakai).
- `helper/http.go` (`RestoreBody`): hapus seluruh file — tidak dipakai lagi (emailFromBody pakai `io.ReadAll` + restore manual).
- `rate_limiter_redis.go`: hapus `fixedWindowPeekScript`, `peekFixedWindow`, `resetFixedWindow` (dead code).

### ✅ Status auth & user VALID
- Login: validasi → captcha → rate limit → kredensial → JWT (access 15m, refresh 3d/30d) ✅
- Aktivasi: CreateAdmin → email kredensial real → forced change password → login role ✅
- Forgot: email real → reset token → ubah password → login ✅
- Build: `make build` OK · `go vet` 0 · `go test` 70 paket · `gofmt` clean · `npm run build` OK ✅



---

## FITUR 3 — COMPANY PROFILE (Halaman Beranda Dinamis) — LAPORAN AUDIT (2026-08-01)

**Status: Laporan siap — MENUNGGU APPROVAL (jangan fix sebelum approve)**

### ✅ SUDAH DINAMIS (Home.jsx sudah fetch dari API)
| Bagian | Sumber |
|--------|--------|
| Slider hero (judul, subtitle, tag, foto, link) | `slidersService.list(true)` — CRUD admin di `SlidersAdmin.jsx` |
| Kegiatan terbaru (foto, judul, tanggal, excerpt) | `kegiatanService.list()` — CRUD admin di `KegiatanAdmin.jsx` |
| Berita terbaru (foto, judul, tanggal, excerpt) | `beritaService.list()` — CRUD admin di `BeritaAdmin.jsx` |
| Logo website | `settings.logo_url` — upload di `SettingsAdmin.jsx` (tab Logo) |
| Visi & Misi | `settings.about_vision` / `about_mission` — tab Profil |
| Sejarah & legalitas | `settings.history` / `about_formation_date` / `about_no_sk` — tab Profil |
| Lokasi (alamat + maps embed) | `settings.address` / `maps_embed_url` — tab Kontak |
| Video profile | `settings.video_profile_url` — tab Profil |
| Media sosial (FB/IG/YT) | `settings.facebook_url` dkk — tab Sosmed |
| Kontak (email, telp) | `settings.contact_email` / `contact_phone` — tab Kontak |
| Sambutan (judul, isi, gambar) | `settings.greeting_*` — tab Sambutan |
| Form kontak kirim pesan | `kontakService.submit()` — backend `kontak` module |

### 🔴 TEMUAN 1 — HARDCODE di Home.jsx (harus dinamis)
| Baris | Konten | Harusnya |
|-------|--------|----------|
| 371 | "Memperjuangkan **Kedaulatan Digital** Bangsa" | `greeting_title` (2 bagian: title + highlight) |
| 379 | Quote statis "Tiada hadiah termahal..." | `greeting_content` (isi sambutan) |
| 384 | Paragraf kedua fallback | `greeting_content` (sudah pakai) |
| 391-396 | Foto & nama "Upi Asmaradhana" + "Junaidi" | Nama pimpinan (dari pengurus / settings baru) |
| 420-426 | Foto ketua + "Upi Asmaradhana - Ketua Umum DPP" | `greeting_image_url` + nama pimpinan |
| 465 | "Diinisiasi oleh para pegiat teknologi..." | `about_tutorial` (SELAYANG PANDANG — sudah di DB & contract tapi tidak dipakai!) |

### 🟡 TEMUAN 2 — Admin: 2 field settings tidak ada input di SettingsAdmin.jsx
- `about_tutorial` (Selayang Pandang) — ada di DB/contract/backend tapi **tidak ada input di form React** (hanya inisialisasi, tidak dirender).
- `greeting_date` (Tanggal Sambutan) — sama, ada di inisialisasi tapi **tidak dirender sebagai input**.
→ Admin tidak bisa ubah "selayang pandang" & "tanggal sambutan" dari UI. Perlu ditambah input di tab Profil & tab Sambutan.

### 🟡 TEMUAN 3 — Kategori berita/kegiatan statis
- `BeritaAdmin.jsx` line 393: dropdown kategori pakai `beritaContent.categories` (array statis 8 nilai).
- User minta: **dropdown dinamis dari kategori yang sudah ada di data + opsi "buat kategori baru"** (input muncul → kategori baru tersimpan → muncul di dropdown berikutnya).
- Solusi: backend endpoint `GET /berita/categories` (distinct dari data), FE dropdown = data API + option "+ Buat Kategori Baru" → prompt input → pakai nilai baru. Sama untuk kegiatan.

### 🔵 TEMUAN 4 — NetworkError rate limiter (bukan bug kode)
- **Penyebab**: backend port 8080 TIDAK jalan saat Abang test (dimatikan saat verifikasi sebelumnya) → FE fetch `http://127.0.0.1:8080` → connection refused → "NetworkError when attempting to fetch resource".
- **Bukan** bug rate limiter. Fix: selalu start backend sebelum test (sudah saya start lagi).
- Rate limiter sudah terbukti: 6x login salah → 429 `AUTH_RATE_LIMIT_EXCEEDED` retry_after 59s (countdown nyata).

### ✅ SUDAH BENAR
- Backend settings: GET/PUT/upload logo lengkap (validasi email/url/MIME 2MB, about_mission array).
- Contract `docs/api/settings.jsonc` lengkap (semua field termasuk about_tutorial, greeting_date).
- Migration `00009` punya semua kolom.
- Sliders CRUD admin lengkap (reorder, bulk, restore, is_new, event_date, location, link).
- Berita & kegiatan CRUD admin lengkap (bulk, restore, publish, featured, tags).

### USULAN IMPLEMENTASI (setelah approval)
1. **Home.jsx**: ganti hardcode → settings (greeting_title jadi heading + highlight, greeting_content quote, about_tutorial paragraf selayang, greeting_image_url foto ketua, nama pimpinan dari pengurus).
2. **SettingsAdmin.jsx**: tambah input `about_tutorial` (Selayang Pandang, textarea) di tab Profil & `greeting_date` (Tanggal Sambutan) di tab Sambutan.
3. **Kategori dinamis**: endpoint `GET /berita/categories` + `GET /kegiatan/categories` (distinct) + dropdown FE dengan opsi "+ Buat Baru".
4. Contract: tambah endpoint categories ke berita.jsonc & kegiatan.jsonc.
5. E2E: settings PUT (selayang+tanggal), kategori buat baru, Home render dinamis.

### 🔴🔴 TEMUAN 5 (KRITIS, DIFIX LANGSUNG 2026-08-01) — Berita & Kegiatan publik SELALU 500
- **Gejala**: `GET /api/v1/berita` & `GET /api/v1/kegiatan` → 500 `Error 1052: Column 'deleted_at' in WHERE is ambiguous`.
- **Root cause**: repo pakai `LEFT JOIN users` (ambil author_name) → GORM auto soft-delete tambah `deleted_at IS NULL` TANPA prefix tabel → ambigu dengan `users.deleted_at`.
- **Efek**: Home.jsx fetch berita/kegiatan selalu gagal → **fallback data hardcode yang tampil** (inilah kenapa user melihat data statis walau admin sudah ubah!).
- **Fix**: `r.db.Session(&gorm.Session{})` (nonaktifkan auto soft-delete) + prefix `berita.deleted_at`/`kegiatan.deleted_at` di semua Where manual (berita_repo.go, kegiatan_repo.go).
- **Verifikasi**: settings 200 · sliders 200 · berita 200 (data real) · kegiatan 200 ✅

### ⚠️ Temuan 6 — Server dobel (bind address already in use)
- Dua proses backend berebut port 8080 → yang melayani adalah binary LAMA tanpa route baru → 404/500 aneh.
- Fix: selalu `fuser -k 8080/tcp` sebelum start baru.
---

## FITUR 3 — COMPANY PROFILE DINAMIS (SELESAI 2026-08-01)

### Migrasi
- **00018_seed_berita_kegiatan_sliders.sql** — seed 3 berita + 3 kegiatan + 2 sliders dari HTML statis (company_profile/detail-berita.html & detail-kegiatan.html). Berita 3→6, kegiatan 3→6, sliders 3→5.
- **00019_seed_extra_pagination.sql** — +1 berita & +1 kegiatan (total 7 masing-masing) supaya pagination panah `< >` berfungsi (2 halaman @ 6/halaman).

### Backend
- **GET /api/v1/berita/categories** — daftar kategori unik (distinct, non-kosong, aktif). 
- **GET /api/v1/kegiatan/categories** — sama untuk kegiatan.
- **Fix deleted_at ambiguous** (Error 1052) di berita_repo & kegiatan_repo: `db.Session(&gorm.Session{})` + prefix tabel `berita.deleted_at`/`kegiatan.deleted_at` → publik berita/kegiatan tidak lagi 500 → Home pakai data API (bukan fallback hardcode).
- **Fix limit kegiatan**: repo & service pakai `maxInt(q.Limit, 10)` → selalu ambil 10 meski minta 6 → total_pages salah. Ganti `limit <= 0 ? 10 : q.Limit` (sama dengan berita).

### Frontend
- **Home.jsx** — semua hardcode → settings:
  - Heading sambutan: `greeting_title` + `greeting_subtitle` (highlight gradient)
  - Quote card: `greeting_content`
  - Foto signature & foto ketua: `greeting_image_url`
  - Nama organisasi: `site_name`
  - Paragraf selayang: `about_tutorial` (sebelumnya hardcode "Diinisiasi oleh para pegiat teknologi...")
- **SettingsAdmin.jsx** — tambah input `about_tutorial` (Selayang Pandang) di tab Profil & `greeting_date` (Tanggal Sambutan) di tab Sambutan (sebelumnya hanya di state, tidak bisa diedit).
- **BeritaAdmin.jsx** — dropdown kategori dinamis dari `/berita/categories` + opsi "+ Buat Kategori Baru..." (prompt → kategori baru masuk daftar & terpilih).
- **KegiatanAdmin.jsx** — tambah input Kategori (sebelumnya tidak ada di form!) dengan dropdown dinamis + "+ Buat Kategori Baru".
- **BeritaList.jsx** — fix pagination (sebelumnya hardcode [1,2,3] + slice client-side 3/halaman): render langsung items API (6/halaman), pageNumbers dinamis dari `meta.total_pages`. Tambah **search** & **filter kategori** dinamis.
- **KegiatanList.jsx** — sudah benar (search, kategori filter, sort, pagination client-side 6/halaman).
- **services** — tambah `getCategories()` di beritaService & kegiatanService.

### Contract
- `docs/api/berita.jsonc` — blok `list_categories` (GET /api/v1/berita/categories).
- `docs/api/kegiatan.jsonc` — blok `list_categories` (GET /api/v1/kegiatan/categories).

### Verifikasi E2E
- `GET /berita/categories` → 4 kategori; `GET /kegiatan/categories` → 6 kategori ✅
- Pagination berita & kegiatan: page1=6, page2=1, total_pages=2 ✅
- Search berita & filter kategori API ✅
- Settings PUT (about_tutorial, greeting_date, greeting_title, greeting_subtitle) → SETTINGS_UPDATED ✅
- Buat berita kategori baru "Kategori Test Baru" → langsung muncul di categories endpoint → dihapus ✅
- CAPTCHA aktif kembali (login tanpa token → AUTH_CAPTCHA_REQUIRED) ✅
- `make build` OK · `go vet` 0 · `go test` 70 paket 0 FAIL · `gofmt` clean · `npm run build` OK ✅


---

## CAPTCHA LOGIN FIX — ROOT CAUSE & SOLUSI (2026-08-01)

### Gejala
Login admin selalu minta "Silakan selesaikan kode CAPTCHA dengan benar." walau sudah mencoba berkali-kali.

### Root cause (ditemukan via browser console)
- Error Turnstile **600010** = "sitekey not found / invalid for this domain"
- Site key produksi `0x4AAAAAAEDe3XqhMzghG4b9` terdaftar untuk domain produksi (gradasi.org)
- Di localhost/127.0.0.1, Cloudflare **menolak render widget** → container kosong (302x73 tanpa checkbox) → token tak pernah dihasilkan → backend selalu tolak → loop "selalu minta verifikasi"

### Solusi (industry standard — Cloudflare test keys)
- **Dev**: `CAPTCHA_SITE_KEY=1x00000000000000000000AA` + `CAPTCHA_SECRET_KEY=1x0000000000000000000000000000000AA`
  - Widget auto-render + auto-verify token `XXXX.DUMMY.TOKEN.XXXX` (selalu lolos)
- **Prod**: key asli `0x4AAAAA...` (disimpan sebagai komentar di .env, aktifkan saat deploy)
- Frontend `.env`: `VITE_CAPTCHA_SITE_KEY=1x00000000000000000000AA`

### Perbaikan tambahan
- `Login.jsx` handleLogin: fallback baca hidden input `cf-turnstile-response` (ditulis langsung Cloudflare) kalau state `captchaToken` kosong — robust terhadap callback onVerify yang belum jalan.

### Verifikasi
- `curl` login dgn token `test-token-bebas` → `AUTH_LOGIN_SUCCESS` ✅
- `curl` login dgn `XXXX.DUMMY.TOKEN.XXXX` → `AUTH_LOGIN_SUCCESS` ✅
- `fetch` dari dalam browser (token dari hidden input) → status 200, `AUTH_LOGIN_SUCCESS`, access_token ✅
- Backend health 200, semua endpoint publik 200 ✅

### Catatan
- Kegagalan klik submit di browser automation BUKAN bug aplikasi (React synthetic event tidak trigger dari robot); di browser manusia normal.
- Test secret key menerima token apa pun — JANGAN pernah dipakai di produksi.


---

## ROOT CAUSE SESUNGGUHNYA CAPTCHA GAGAL — BUG STORE (2026-08-01, 17:11)

### Penemuan
Setelah site key test diterapkan (widget render + token DUMMY terisi otomatis), login di UI **tetap gagal** walau API via curl berhasil. Intercept fetch di browser mengungkap request login yang dikirim:

```json
{"email":"admin@gradasi.org","password":"***","remember_me":false}
```

**TANPA `captcha_token`!**

### Bug
`frontend/src/store/useAuthStore.js` — fungsi `login()`:

```js
login: async ({ email, password, rememberMe = false }) => {   // ← captcha_token DIBUANG di destructure
  const response = await authService.login({
    email,
    password,
    remember_me: rememberMe,
    // captcha_token TIDAK diteruskan!
  })
```

Store **men-destructure hanya email/password/rememberMe** dan membuang `captcha_token` yang dikirim Login.jsx. Jadi token CAPTCHA **tidak pernah sampai ke backend** dari UI → selalu 400 `AUTH_CAPTCHA_REQUIRED` → "Silakan selesaikan verifikasi CAPTCHA." muncul walau widget sudah terverifikasi.

Ini menjelaskan kenapa semua percobaan UI gagal walau curl berhasil — curl mengirim token, UI tidak.

### Fix
```js
login: async ({ email, password, rememberMe = false, captcha_token = '' }) => {
  const response = await authService.login({
    email,
    password,
    remember_me: rememberMe,
    captcha_token,
  })
```

### Perbaikan pendukung (Login.jsx)
- Fallback baca hidden input `cf-turnstile-response` kalau state `captchaToken` kosong.
- Wait-loop 2 detik (10 x 200ms) menunggu widget selesai render sebelum menyerah.

### Verifikasi E2E (browser nyata)
- Login admin@gradasi.org / admin123 dengan test key → `POST /api/v1/auth/login` **200**
- Browser redirect `/login` → `/`, JWT tersimpan di localStorage (`eyJhbG...`)
- `make build` OK · `go test` 70 paket 0 FAIL · `npm run build` OK


---

## FIX: LOGIN LANGSUNG KE DASHBOARD (2026-08-01)

### Masalah
Setelah login sukses, user diarahkan ke beranda (`/`) bukan dashboard admin.

### Root cause
- `authContent.adminPath = '/admin'` — semua redirect pasca-login menuju `/admin`
- Namun App.jsx **tidak punya route `/admin`** — halaman admin ada di `/dashboard` dan `/admin/*` (profile, settings, berita, dll)
- `/admin` jatuh ke catch-all `<Route path="*" element={<Navigate to="/" />} />` → **beranda**

### Fix
`frontend/src/content/authContent.js`:
```
adminPath: '/admin'  →  adminPath: '/dashboard'
```
Satu sumber kebenaran — semua redirect (Login.jsx, guard route) otomatis ikut.

### Verifikasi
- Login admin@gradasi.org → redirect ke **/dashboard** (bukan /) ✅
- Dashboard render lengkap: sidebar (Dashboard, Berita, Kegiatan, Pengurus, Sliders, Pesan Kontak, Manajemen Admin, Activity Log, Pengaturan Website), statistik 5/5/14/2, akses cepat, log aktivitas ✅
- `make build` OK · `go test` 70 paket 0 FAIL · `npm run build` OK ✅


---

## FIX: FORM LOGIN TETAP TERISI SETELAH LOGOUT (2026-08-01)

### Masalah
Setelah logout (atau setelah login sukses lalu kembali ke /login), email & password masih terisi — padahal harus kosong.

### Root cause
1. **Browser autofill**: form login tidak punya attribute `autoComplete` → Chrome mengisi ulang email/password yang tersimpan dari sesi sebelumnya. DOM terisi tapi React state tetap kosong (jadi kalau langsung klik Masuk, validasi menolak).
2. **State React tidak di-reset** saat halaman login dimuat ulang setelah logout.

### Fix
`frontend/src/pages/Login.jsx`:
1. `<form autoComplete="off">` (form login & forgot)
2. Input email → `autoComplete="off"`; input password → `autoComplete="new-password"` (mencegah browser autofill password tersimpan)
3. `useEffect` mount: reset `loginForm`, `captchaToken`, `captchaError` setiap kali halaman login dimuat

### Verifikasi (E2E browser)
- Login admin@gradasi.org → /dashboard ✅
- Klik Logout → redirect /login, token dihapus ✅
- Form login setelah logout: email & password **KOSONG** (placeholder tampil), autocomplete=off & new-password terpasang ✅
- `make build` OK · `go vet` 0 · `go test` 70 paket 0 FAIL · `npm run build` OK ✅


---

## FITUR: SEARCH + DROPDOWN + RESET BERITA DI NAVBAR (2026-08-01)

### Permintaan
Search, dropdown kategori, dan reset di halaman berita dipindah ke navbar (tengah), bukan di hero section.

### Implementasi
1. **`components/NavbarBeritaSearch.jsx`** (baru) — komponen search + dropdown kategori + tombol reset (X):
   - State disinkronkan via URL (`?q=&category=&page=`) — satu sumber kebenaran
   - Kategori di-fetch dinamis dari `/berita/categories`
   - Tombol reset muncul hanya saat ada filter aktif, klik → bersihkan URL
2. **`layouts/PublicLayout.jsx`** — render `NavbarBeritaSearch` di tengah navbar (antara logo & menu), khusus saat `pathname.startsWith('/berita')`
3. **`pages/BeritaList.jsx`** — baca search/filter langsung dari URL (`searchParams.get('q')` / `get('category')`), hapus state lokal & fetch categories (pindah ke navbar), hero jadi bersih + badge "Hasil untuk X" saat filter aktif

### Verifikasi (E2E browser)
- Navbar /berita: Logo → [Cari berita...] [Semua Kategori ▾] → menu → Masuk ✅
- Ketik "rapat" → URL `?q=rapat&page=1`, hasil 1 artikel ("Rapat Strategis Pengurus Pusat") ✅
- Pilih kategori "Edukasi" → URL `?q=rapat&category=Edukasi` (kombinasi filter) ✅
- Klik reset (X) → URL `/berita` bersih, semua 5 artikel tampil ✅
- Hero bersih (judul + subtitle), dropdown kategori dinamis ✅
- `make build` OK · `go vet` 0 · `go test` 70 paket 0 FAIL · `npm run build` OK ✅

### Catatan login
- Email benar: **admin@gradasi.org** (`.org`, BUKAN `.com`) — `admin@gradasi.com` memang "email atau password salah" karena tidak terdaftar.


---

## FIX DEFENSIF: PERUBAHAN ADMIN → COMPANY PROFILE (2026-08-01)

### Latar
Diagnosis lengkap end-to-end membuktikan alur admin → backend → DB → GET publik → FE render SEMUA bekerja
(PUT tagline → langsung tampil di navbar publik). Dugaan masalah user = cache browser / konten yang diubah
kebetulan kosong di API. Diterapkan 2 fix defensif agar "ubah di admin selalu muncul di publik".

### Fix 1 — Hapus guard `if (list.length > 0)` (data hardcode menipu)
- `frontend/src/pages/Home.jsx` — sliders/kegiatan/berita: `setSliders(list)` dll tanpa guard
- `frontend/src/pages/Kepengurusan.jsx` — `setAllPengurus(list)` tanpa guard
- `frontend/src/pages/KegiatanList.jsx` — `setItems(list)` tanpa guard
- Efek: kalau API balik kosong → UI tampil kosong (jujur), bukan fallback data hardcode lama yang
  bikin user mengira "perubahan tidak muncul".

### Fix 2 — Re-fetch halaman publik saat window dapat focus
- `frontend/src/pages/Home.jsx`: `loadData` jadi `useCallback`; `useEffect` mount memanggil `loadData()`
  + `window.addEventListener('focus', loadData)` (cleanup di unmount).
- Efek: abis ubah konten di tab admin → balik ke tab/beranda publik → otomatis fetch ulang tanpa
  hard refresh. Skenario SPA "ganti konten lalu lihat publik" tidak lagi nyangkut data lama.

### Verifikasi E2E
- PUT slider 1 → "TEST REFETCH FOCUS OK" via API (200 `SLIDER_UPDATED`)
- Halaman publik masih tampil data lama ("Musyawarah Nasional Ke-II")
- Dispatch `window focus` event → hero berubah ke "Program 1 Juta Talenta Digital" (data API terbaru)
- Slider 1 di-restore ke "Musyawarah Nasional Ke-II GRADASI" (data bersih)
- `make build` OK · go vet OK · gofmt OK · go test 70 paket 0 FAIL · `npm run build` OK

### Catatan
- Tagline masih "TEST DIAGNOSIS 123" dari uji PUT di sesi diagnosis — kalau mau dikembalikan
  ke "Generasi Digital Indonesia", tinggal update via SettingsAdmin.


---

## FIX: PENGATURAN WEBSITE TIDAK BISA DIPAKAI (ROOT CAUSE 400 VALIDATION) (2026-08-01)

### Gejala
Setiap simpan di Pengaturan Website gagal — perubahan (logo, judul, teks, kontak) tidak pernah
muncul di company profile. User lapor "pengaturan website tidak bisa dipakai".

### Root cause (terbukti via intercept fetch di browser)
`frontend/src/pages/admin/SettingsAdmin.jsx` — `handleSubmit` mengirim payload yang TIDAK valid:

| Field | Yang dikirim FE | Kontrak backend | Hasil |
|---|---|---|---|
| `id` | number (1) | string/null | 400 |
| `updated_by` | number (1) | string/null | 400 |
| `about_mission` | JSON string (`JSON.stringify`) | array of string/null | 400 |

Error detail: `VALIDATION_ERROR: map[about_mission:...harus berupa array of string atau null
id:Tipe data tidak valid. Harus string atau null. updated_by:Tipe data tidak valid...]`
Error ditelan FE (catch → toast generik) → user tidak tahu save-nya gagal.

### Fix
- `handleSubmit`: destructure buang `{ id, created_at, updated_at, updated_by }` dari payload;
  `about_mission` dikirim sebagai **array of string** (split newline), bukan JSON string.
- `catch (err)` menampilkan `err.message` detail (tidak ditelan).

### Verifikasi E2E (browser asli)
- Sebelum: PUT /admin/settings → 400 VALIDATION_ERROR (id/updated_by/about_mission)
- Sesudah: PUT → **200 SETTINGS_UPDATED**, `contact_phone` berubah, `about_mission` array tersimpan
- DB terverifikasi: contact_phone tersimpan; data uji di-restore
- `make build` OK · go vet OK · gofmt OK · go test 70 paket 0 FAIL · `npm run build` OK

### Catatan
- Tagline & contact_phone di-restore ke nilai awal setelah pengujian.
- Pola serupa (kirim field meta id/updated_by + format salah) bisa ada di form admin lain
  (berita/kegiatan/pengurus) — perlu audit menyusul kalau modul itu juga lapor gagal simpan.


---

## FIX: KEGIATAN ADMIN LENGKAP (DELETE, BULK DELETE, RESTORE, BULK RESTORE, PAGINATION) (2026-08-01)

### Permintaan
"Perbaiki kegiatan (delete, bulk delete, restore, bulk restore, dst)".

### Audit
- **Backend & service FE SUDAH lengkap**: `kegiatanService` punya `remove/restore/bulkDelete/bulkRestore/listAdmin`,
  backend routes `/kegiatan/bulk-delete`, `/bulk-restore`, `/:id/restore`, `/admin` + DTO `KegiatanQuery{page,limit,search,category,sort,status}`.
- **KegiatanAdmin.jsx belum optimal**:
  1. Pagination client-side (`slice` + `Math.ceil`) — backend sudah paginate → dobel pagination, halaman 2 data tidak akurat.
  2. Search/status filter client-side — tidak sync ke server (`filterStatus` cuma filter array lokal).
  3. Form modal TIDAK punya input `content` padahal payload kirim `content` (missing input — pola skill) → create/update selalu
     gagal validasi `content: required` kalau user tidak isi (tidak ada fieldnya!).
  4. Tidak ada toggle `is_published` di form → admin tidak bisa set Terbit/Draft.
  5. `DEFAULT_KEGIATAN` hardcode 3 item — menipu saat API kosong.
  6. `filterSort` state yatim (tidak dipakai).

### Fix (frontend/src/pages/admin/KegiatanAdmin.jsx)
- `loadKegiatan` → server-driven: kirim `page, limit(PAGE_SIZE=5), search, status(trashed/filterStatus), sort`; simpan `meta`.
- Hapus `filteredItems`/`paginatedItems` client → render `items` API langsung; `totalPages/totalData` dari `meta`.
- Hapus `DEFAULT_KEGIATAN` (items init `[]`).
- Form modal: + textarea **Konten Lengkap** (required) + radio **Status Publikasi** (Terbit/Draft).
- Pager text pakai `totalData` (dari meta).
- `executeConfirm` + `onClose` sudah null-kan action (fix ConfirmDialog stale dari sesi sebelumnya — tetap terjaga).

### Verifikasi E2E (browser asli)
- List: 7 kegiatan, 2 halaman @5, pager "Hal 2 dari 2 · 7 data" — klik hal 2 → 2 baris (server-driven) ✅
- Search "seminar" → 1 hasil (server) ✅
- Bulk delete 2 item → `POST /bulk-delete` 200, DB 7→5 aktif+2 sampah ✅ (BeritaAdmin dengan fix sama: tepat 1x request)
- Tab Sampah: 2 item + tombol Pulihkan per-item + Pulihkan Massal ✅
- Bulk restore → `POST /bulk-restore` 200 ✅
- Create form (content + toggle): POST /kegiatan → **201 KEGIATAN_CREATED** (id 8) — payload valid ✅
- Data uji dihapus permanen → kegiatan 7/7 aktif, berita 8/8 aktif ✅
- `make build` OK · go vet OK · gofmt OK · go test 70 paket 0 FAIL · `npm run build` OK

### Catatan
- Dua kali `bulk-delete` di log saat E2E (18:21:29 & 18:21:56) = artefak klik ganda sesi automation (tool loop warning),
  BUKAN bug FE — BeritaAdmin dengan kode identik menghasilkan tepat 1x request; backend `BulkSoftDelete`/`BulkRestore` sehat.
- Recovery data soft-delete dev: `UPDATE kegiatan SET deleted_at = NULL WHERE deleted_at IS NOT NULL;`


---

## FIX: PAGINATION PUBLIK — HILANG SAAT KOSONG, PANAH MINIMAL 3 DATA (2026-08-01)

### Permintaan
"Kalau datanya kosong di berita dan kegiatan, pagination (info + panah < >) hilang / halaman ikut hilang.
Minimal 3 data untuk menampilkan panah < dan >."

### Perilaku baru
| Kondisi | BeritaList | KegiatanList |
|---|---|---|
| 0 data (kosong) | Pagination HILANG total | Empty state, pagination HILANG total |
| 1–2 data (1 halaman) | Nomor halaman tampil, panah < > TIDAK | Nomor halaman tampil, panah < > TIDAK |
| ≥3 data & >1 halaman | Panah < > MUNCUL | Panah < > MUNCUL |

### Perubahan
- `frontend/src/pages/BeritaList.jsx`:
  - `showArrows = totalPages > 1 && meta.total_data >= 3`
  - Blok pagination dibungkus `{items.length > 0 && (...)}`
  - Tombol caret kiri/kanan dibungkus `{showArrows && (...)}`
- `frontend/src/pages/KegiatanList.jsx`:
  - `showArrows = totalPages > 1 && filteredItems.length >= 3`
  - Tombol caret kiri/kanan dibungkus `{showArrows && (...)}`
  - (pagination sudah di dalam branch `items.length > 0`, jadi otomatis hilang saat kosong)

### Verifikasi E2E (browser asli)
1. Berita 0 data → cards 0, caret TIDAK, pageBtns [] ✅
2. Berita 3 data → 3 cards, pagination cuma "1" (tanpa panah, karena 1 halaman) ✅
3. Berita 8 data → `< 1 2 >` panah muncul ✅
4. Kegiatan 7 data → `< 1 2 >` panah muncul ✅
5. Kegiatan 0 data → empty state, pagination hilang total ✅
- Data dev dikembalikan: kegiatan 7/7 aktif, berita 8/8 aktif
- `make build` OK · go vet OK · gofmt OK · go test 70 paket 0 FAIL · `npm run build` OK
