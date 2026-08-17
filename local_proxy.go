package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

// Ini adalah alat bantu untuk simulasi NGINX di lokal laptop kamu.
// Nanti di VPS/Production, script ini TIDAK DIPAKAI. Di VPS kamu murni pakai config NGINX asli.

func main() {
	// Target Frontend (Vite) jalan di 5173
	frontendURL, _ := url.Parse("http://localhost:5173")
	frontendProxy := httputil.NewSingleHostReverseProxy(frontendURL)

	// Target Backend (Go) jalan di 8080
	backendURL, _ := url.Parse("http://localhost:8080")
	backendProxy := httputil.NewSingleHostReverseProxy(backendURL)

	// Deteksi Bot Sosial Media (Sama persis dengan settingan Nginx VPS)
	botRegex := regexp.MustCompile(`(?i)(WhatsApp|facebookexternalhit|Twitterbot|LinkedInBot|TelegramBot)`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")

		// Deteksi apakah manusia atau bot
		isBot := botRegex.MatchString(userAgent)
		
		// Deteksi URL yang diakses
		isKegiatan := len(r.URL.Path) >= 9 && r.URL.Path[:9] == "/kegiatan"
		isBerita := len(r.URL.Path) >= 7 && r.URL.Path[:7] == "/berita"

		// 1. RULE SEO BOT PROXY: Kalau bot mengakses URL berita/kegiatan
		if isBot && (isKegiatan || isBerita) {
			log.Printf("[NGINX SIMULATOR] Bot %s membagikan %s -> Dibelokkan ke Backend SEO", userAgent, r.URL.Path)
			r.URL.Path = "/api/v1/seo" + r.URL.Path
			backendProxy.ServeHTTP(w, r)
			return
		}

		// 2. RULE API: Semua request API diarahkan ke backend 8080
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			backendProxy.ServeHTTP(w, r)
			return
		}

		// 3. RULE UPLOADS: Ambil gambar langsung dari backend
		if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/uploads" {
			backendProxy.ServeHTTP(w, r)
			return
		}
		
		// 4. RULE DEFAULT: Selain itu semua (manusia), diarahkan ke Frontend React 5173
		frontendProxy.ServeHTTP(w, r)
	})

	log.Println("=====================================================")
	log.Println("✅ NGINX Simulator (Local Proxy) Berjalan di Port 8888")
	log.Println("Arahkan Ngrok kamu ke 8888, bukan 5173!")
	log.Println("=====================================================")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
