# API Contract — Dynamic Content Management (Site Settings)

> Modul: **settings** — konten branding/copy website (logo, nama aplikasi, tagline,
> kontak, sosial media, section teks) yang dikelola admin dan dirender di
> company profile publik.
>
> Base URL: `/api/v1` · Format: JSON · Prefix route: `/settings` (publik) dan `/admin/settings` (admin)
> Source of truth tambahan: `docs/api/settings.jsonc` (versi JSONC, sama isinya).

---

## Ringkasan Endpoint

| # | Method | Endpoint | Auth | Fungsi |
|---|--------|----------|------|--------|
| 1 | GET  | `/api/v1/settings` | Publik (tanpa auth) | Ambil semua settings untuk render company profile |
| 2 | PUT  | `/api/v1/admin/settings` | Bearer JWT + Role `super_admin`/`admin` | Update sebagian/keseluruhan field settings |
| 3 | POST | `/api/v1/admin/settings/logo` | Bearer JWT + Role `super_admin`/`admin` | Upload logo (multipart), simpan ke disk lokal |

Security bersama (endpoint 2 & 3): middleware `AuthMiddleware` (JWT),
`RoleMiddleware("super_admin","admin")`, rate limit `RateLimitPerUser("settings-admin", 30)`.

---

## 1. GET /api/v1/settings

Ambil seluruh konfigurasi website. Data publik — bisa diakses tanpa login.
Response inilah yang di-`binding` ke UI statis company profile (logo, judul,
tagline, kontak, sosial media, teks section).

### Response 200

```json
{
  "success": true,
  "code": "SETTINGS_RETRIEVED",
  "message": "Konfigurasi website berhasil diambil",
  "data": {
    "id": 1,
    "site_name": "DPP GRADASI",
    "tagline": "Generasi Digital Indonesia",
    "logo_path": "/uploads/settings/1737187847.png",
    "contact_email": "dpp@gradasi.org",
    "contact_phone": "+6281234567890",
    "address": "Office Park OL3-IZA The Bellagio Mall, Mega Kuningan, Jakarta Selatan",
    "maps_embed_url": "https://www.google.com/maps/embed?pb=...",
    "facebook_url": "https://www.facebook.com/gradasiofficial.id",
    "instagram_url": "https://www.instagram.com/dppgradasi",
    "youtube_url": "https://www.youtube.com/channel/UCwdjB4LkqcF4Kw5-PoyOb5A",
    "video_profile_path": "https://www.youtube.com/embed/dQw4w9WgXcQ",
    "history": "Perkumpulan Generasi Digital Indonesia (GRADASI) didirikan pada 4 Februari 2019...",
    "about_tutorial": "Pengesahan Badan Hukum Kemenkumham RI.",
    "about_formation_date": "4 Februari 2019",
    "about_no_sk": "AHU-0000151.AH.01.07.2019",
    "about_vision": "Mewujudkan masyarakat Indonesia yang cerdas, kreatif, dan berdaulat di era digital.",
    "about_mission": ["Membangun ekosistem literasi digital yang inklusif di seluruh daerah Indonesia."],
    "greeting_title": "Tahun Baru 2026",
    "greeting_subtitle": "Resolusi & Harapan",
    "greeting_date": "11 Februari 2026",
    "greeting_content": "Memasuki tahun 2026, GRADASI menetapkan pilar utama perjuangan...",
    "greeting_image_path": "https://gradasi.org/uploads/img/event-terkini/1767154211.jpg",
    "created_at": "2026-07-30T08:00:00Z",
    "updated_at": "2026-07-30T08:00:00Z",
    "updated_by": 1
  }
}
```

### Catatan field
- Semua field opsional untuk dirender — FE wajib punya fallback bila kosong.
- `logo_path`, `greeting_image_path`, dan URL upload lokal memakai path relatif
  (`/uploads/...`) yang di-serve Gin; URL eksternal penuh (`https://...`) juga
  didukung. FE cukup prepend base URL bila path relatif.
- `updated_by` = ID user admin terakhir yang mengubah settings (nullable).

---

## 2. PUT /api/v1/admin/settings

Update sebagian/keseluruhan field. Body berupa **partial object** — key yang
tidak dikirim tidak diubah. Key yang tidak dikenal → 422.

### Request

```json
{
  "site_name": "DPP GRADASI",
  "tagline": "Generasi Digital Indonesia",
  "contact_email": "dpp@gradasi.org",
  "contact_phone": "+6281234567890",
  "about_vision": "Mewujudkan masyarakat Indonesia yang cerdas, kreatif, dan berdaulat di era digital.",
  "about_mission": ["Membangun ekosistem literasi digital yang inklusif di seluruh daerah Indonesia."]
}
```

### Validasi per field
- Key harus dikenal (snake_case sesuai DTO). Key tak dikenal → 422 `VALIDATION_ERROR`, errors berisi `{key: "Field tidak dikenal"}`.
- Field berakhiran `_email` → wajib format email valid.
- Field berakhiran `_url` → wajib diawali `http://` atau `https://`.
- `about_mission` → array of string **atau** `null` (disimpan sebagai JSON string di DB). Tipe lain → 422.
- Field lain → string (atau `null` → di-set kosong). Tipe lain → 422.
- `updated_by` di-set otomatis dari JWT (tidak boleh dikirim client).

### Response 200

```json
{
  "success": true,
  "code": "SETTINGS_UPDATED",
  "message": "Konfigurasi website berhasil diperbarui",
  "data": { "site_name": "DPP GRADASI", "...": "semua field seperti GET /settings" }
}
```

### Response 422

```json
{
  "success": false,
  "code": "VALIDATION_ERROR",
  "message": "Validasi gagal.",
  "errors": [
    { "field": "contact_email", "tag": "invalid", "message": "Format email tidak valid" },
    { "field": "siteName", "tag": "invalid", "message": "Field tidak dikenal" }
  ]
}
```

### Response 401 / 403 / 429
- 401: `{"success":false,"code":"AUTH_TOKEN_INVALID_OR_EXPIRED","message":"..."}` (token invalid/expired)
- 403: `{"success":false,"code":"FORBIDDEN","message":"..."}` (role bukan super_admin/admin)
- 429: `{"success":false,"code":"AUTH_RATE_LIMIT_EXCEEDED","message":"...","retry_after":60}`

---

## 3. POST /api/v1/admin/settings/logo

Upload file logo baru (multipart/form-data). Validasi: wajib file, size maks
**2 MB**, MIME type **PNG/JPEG/WEBP**. File disimpan ke `public/uploads/settings/`
dengan nama unik (`<unixnano><ext>`), lalu `logo_path` di tabel settings
di-update ke path relatif `/uploads/settings/<nama>`.

### Request (multipart/form-data)

```
logo: (binary) — wajib, maks 2MB, image/png | image/jpeg | image/webp
```

### Response 200

```json
{
  "success": true,
  "code": "SETTINGS_LOGO_UPLOADED",
  "message": "Logo berhasil diunggah",
  "data": { "site_name": "DPP GRADASI", "logo_path": "/uploads/settings/1737187847123456789.png", "...": "seluruh field seperti GET /settings" }
}
```
> Catatan: backend mengembalikan seluruh objek settings (sama seperti GET /settings), bukan hanya `logo_path`.

### Response 400 (file salah)

```json
{
  "success": false,
  "code": "INVALID_LOGO",
  "message": "File logo tidak valid. Maksimal 2MB dengan format PNG, JPG, atau WEBP.",
  "errors": null
}
```

### Response 422 (logo tidak dikirim)

```json
{
  "success": false,
  "code": "VALIDATION_ERROR",
  "message": "Validasi gagal.",
  "errors": [ { "field": "logo", "tag": "required", "message": "File logo wajib diunggah." } ]
}
```

---

## Alur Sinkronisasi FE-BE

1. **PublicLayout** (header/footer) + **Home** fetch `GET /settings` sekali lewat `SettingsContext` → binding logo, site_name, tagline, sosial media, alamat, video, dll.
2. **SettingsAdmin** fetch `GET /settings` → isi form; submit `PUT /admin/settings` (partial); upload logo `POST /admin/settings/logo` (FormData).
3. Setelah update sukses, FE invalidate cache context → company profile langsung tampil data baru tanpa refresh manual.
