package db

import (
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
)

type Connection interface {
	GetDatabaseConnection() *gorm.DB
}

type connection struct {
	db *gorm.DB
}

var (
	db   *gorm.DB
	once = &sync.Once{}
)

func (c *connection) GetDatabaseConnection() *gorm.DB {
	return c.db
}

func NewConnection(dns string) (Connection, error) {

	once.Do(func() {
		var err error
		db, err = gorm.Open(os.Getenv("DB_CONNECTION"), dns)
		if err != nil {
			panic(err)
		}

		db.AutoMigrate(&entities.User{}, &entities.Erce{}, &entities.Erfce{}, &entities.Warning{})
	})

	return &connection{db: db}, nil
}
