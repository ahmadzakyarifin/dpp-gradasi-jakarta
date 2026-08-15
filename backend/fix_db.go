package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Mama14081980@tcp(localhost:3306)/dpp_new?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("DB Connect Error:", err)
		return
	}

	var pengurusList []model.Pengurus
	db.Find(&pengurusList)

	re := regexp.MustCompile(`(?i)(?:\s|^)(\d+\.\s)`)

	for _, p := range pengurusList {
		changed := false

		if p.Sertifikasi != nil && *p.Sertifikasi != "" {
			newSert := re.ReplaceAllString(*p.Sertifikasi, "\n\n$1")
			newSert = strings.TrimSpace(newSert)
			if newSert != *p.Sertifikasi {
				p.Sertifikasi = &newSert
				changed = true
			}
		}

		if p.Pendidikan != nil && *p.Pendidikan != "" {
			newPend := re.ReplaceAllString(*p.Pendidikan, "\n\n$1")
			newPend = strings.TrimSpace(newPend)
			if newPend != *p.Pendidikan {
				p.Pendidikan = &newPend
				changed = true
			}
		}

		if changed {
			err = db.Save(&p).Error
			if err != nil {
				fmt.Println("Error updating:", p.ID, err)
			} else {
				fmt.Println("Updated pengurus:", p.ID, p.Name)
			}
		}
	}
	fmt.Println("Done")
}
