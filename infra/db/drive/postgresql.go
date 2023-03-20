package drive

import (
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type IPostgres interface {
	Connect() string
}

type Postgres struct {
	IDrive
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Connect() string {
	
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_NAME := os.Getenv("DB_NAME")
	SSL_MODE := os.Getenv("SSL_MODE")

	return "host=" + DB_HOST + " port=" + DB_PORT + " user=" + DB_USERNAME + " dbname=" + DB_NAME + " password=" + DB_PASSWORD + " sslmode=" + SSL_MODE

}
