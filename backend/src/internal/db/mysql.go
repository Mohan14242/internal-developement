package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitMySQL() {
	log.Println("🗄️ Initializing MySQL connection")

	// Read env vars (DO NOT log password)
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")

	log.Printf(
		"📡 MySQL config → host=%s port=%s user=%s db=%s",
		dbHost, dbPort, dbUser, dbName,
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		os.Getenv("DB_PASSWORD"), // never log this
		dbHost,
		dbPort,
		dbName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println("❌ Failed to open MySQL connection:", err)
		panic(err)
	}

	log.Println("🔌 MySQL connection opened, pinging database...")

	if err = DB.Ping(); err != nil {
		log.Println("❌ MySQL ping failed:", err)
		panic(err)
	}

	log.Println("✅ MySQL connection established successfully")
}