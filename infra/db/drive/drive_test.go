package drive_test

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
	"github.com/risk-place-angola/backend-risk-place/infra/db"
	"github.com/risk-place-angola/backend-risk-place/infra/db/drive"
	"github.com/risk-place-angola/backend-risk-place/infra/db/drive/postgres"
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../../../.env.example")

	if err != nil {
		log.Fatalf("Error loading .env files")
	}
}

func TestDrivePostgres(t *testing.T) {

	os.Setenv("DB_HOST", "localhost")

	postgres := postgres.NewPostgres()
	driver := drive.Drive(postgres)

	if driver == "" {
		t.Error("driver is empty")
	}

	dns := postgres.Connect()

	if dns == "" {
		t.Error("dns is empty")
	}

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
