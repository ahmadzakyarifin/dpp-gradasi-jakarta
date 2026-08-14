// Package email berisi template email HTML untuk DPP Gradasi.
// Template dirender dengan html/template dan disisipkan ke binary lewat embed
// sehingga tidak bergantung pada filesystem saat runtime.
package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed *.html
var templateFS embed.FS

// Render memuat template bernama name (mis. "forgot_password.html"),
// mem-parsingnya dengan data, dan mengembalikan HTML body yang siap dikirim.
//
// Field yang umum dipakai template:
//   - Name    : nama penerima (string)
//   - URL     : tautan aksi (string)
//   - Expired : masa berlaku tautan dalam menit (number)
func Render(name string, data any) (string, error) {
	tmpl, err := template.ParseFS(templateFS, name)
	if err != nil {
		return "", fmt.Errorf("gagal memuat template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("gagal merender template %s: %w", name, err)
	}

	return buf.String(), nil
}
