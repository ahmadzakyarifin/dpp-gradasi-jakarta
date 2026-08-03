package entity

import "time"

// PesanKontak adalah representasi domain dari tabel pesan_kontak.
type PesanKontak struct {
	ID           uint
	Nama         string
	Email        string
	Subjek       string
	Pesan        string
	IsRead       bool
	ResponseNote *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
