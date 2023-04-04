package postgres

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/risk-place-angola/backend-risk-place/infra/db"
	"github.com/risk-place-angola/backend-risk-place/infra/db/drive"
)

func ConnectionPostgres() (*gorm.DB, error) {
	postgres := NewPostgres()
	driver := drive.Drive(postgres)
	if driver == "" {
		return nil, errors.New("driver is empty")
	}

	dns := postgres.Connect()
	if dns == "" {
		return nil, errors.New("dns is empty")
	}

	db, err := db.NewConnection(dns)
	if err != nil {
		return nil, err
	}
	dbConn := db.GetDatabaseConnection()
	if dbConn == nil {
		return nil, errors.New("dbConn is nil")
	}

	return dbConn, nil
}
