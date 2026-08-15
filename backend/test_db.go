package main

import (
	"fmt"
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
	db.Where("level = ?", "Pengurus Kab/Kota").Find(&pengurusList)
	for _, p := range pengurusList {
		fmt.Printf("ID: %d, Nama: %s, Kab: %s\n", p.ID, p.Name, p.Kabupaten)
	}
}
