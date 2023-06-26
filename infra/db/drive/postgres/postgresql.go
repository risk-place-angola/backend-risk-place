package postgres

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/risk-place-angola/backend-risk-place/infra/db/drive"
	"github.com/risk-place-angola/backend-risk-place/util"
)

type IPostgres interface {
	Connect() string
}

type Postgres struct {
	drive.IDrive
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Connect() string {
	env := util.LoadEnv(".env")

	return "host=" + env.DBHOST + " port=" + env.DBPORT + " user=" + env.DBUSER + " dbname=" + env.DBNAME + " password=" + env.DBPASS + " sslmode=" + env.SSLMODE
}
