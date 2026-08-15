package seeder

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedUser(db *sql.DB) {
	username := "haritsnb"
	password := "secretbanget"

	var exists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	err := db.QueryRow(queryCheck, username).Scan(&exists)
	if err != nil {
		log.Printf("Gagal cek seeder user: %v\n", err)
		return
	}

	if !exists {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		queryInsert := `INSERT INTO users (username, password, created_by, modified_by) VALUES ($1, $2, $3, $4)`
		_, err := db.Exec(queryInsert, username, string(hashedPassword), "seeder", "seeder")
		if err != nil {
			log.Printf("Gagal seeding user: %v\n", err)
			return
		}
		fmt.Println("Seeding user berhasil: haritsnb / secretbanget")
	}
}
