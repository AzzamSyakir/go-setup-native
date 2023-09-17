package migration

import (
	"database/sql"
	"log"
)

// Migrate digunakan untuk menjalankan migrasi tabel.
func Migrate(db *sql.DB) {
	// SQL statement untuk memeriksa apakah tabel users sudah ada
	checkTableSQL := `
        SELECT 1 FROM users LIMIT 1
    `

	// Menjalankan perintah SQL untuk memeriksa apakah tabel sudah ada
	var exists int
	err := db.QueryRow(checkTableSQL).Scan(&exists)
	if err == nil {
		// Jika tabel sudah ada, maka tidak perlu melakukan migrasi
		log.Println("Tabel sudah di migrasi")
		return
	}

	// SQL statement untuk membuat tabel users
	createTableSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL,
            email VARCHAR(255) NOT NULL UNIQUE
        )
    `

	// Menjalankan perintah SQL untuk membuat tabel
	_, err = db.Exec(createTableSQL)
	if err != nil {
		// Menangani kesalahan jika terjadi kesalahan saat migrasi
		log.Fatal(err)
		return
	}

	// Pesan sukses jika migrasi berhasil
	log.Println("Migrasi tabel users berhasil")
}
