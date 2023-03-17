package db_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/risk-place-angola/backend-risk-place/infra/db"
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../../.env.example")

	if err != nil {
		log.Fatalf("Error loading .env files")
	}
}

func TestConnectionPostgres(t *testing.T) {

	host := "localhost"
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	dns := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, db_name, password)

	db, err := db.NewConnection(dns)
	if err != nil {
		t.Error(err)
	}
	if db == nil {
		t.Error("db is nil")
	}

	dbConn := db.GetDatabaseConnection()
	if dbConn == nil {
		t.Error("dbConn is nil")
	}

	defer dbConn.Close()

}
