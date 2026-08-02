// Package whatsapp berisi template pesan WhatsApp untuk SchoolPay.
// Template dirender dengan text/template dan disisipkan ke binary lewat embed.
package whatsapp

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed *.txt
var templateFS embed.FS

// Render memuat template bernama name (mis. "activation.txt"),
// mem-parsingnya dengan data, dan mengembalikan body teks yang siap dikirim.
//
// Field yang umum dipakai template:
//   - Name  : nama penerima (string)
//   - URL   : tautan aksi (string)
//   - OTP   : kode verifikasi (string)
//   - Value : nilai nominal (string)
//   - Due   : tanggal jatuh tempo (string)
func Render(name string, data any) (string, error) {
	tmpl, err := template.ParseFS(templateFS, name)
	if err != nil {
		return "", fmt.Errorf("whatsapp: gagal memuat template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("whatsapp: gagal merender template %s: %w", name, err)
	}

	return buf.String(), nil
}
