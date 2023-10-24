package database

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbCredentials struct {
	client   string
	user     string
	password string
	port     string
	host     string
	database string
	ssl      string
}

func (db *dbCredentials) formatStr() string {
	return "user=" + db.user + " host=" + db.host + " port=" + db.port + " password=" + db.password + " dbname=" + db.database + " sslmode=" + db.ssl
}

var Db *gorm.DB

func Connectdb() {
	sql.Drivers()
	creds := dbCredentials{
		client: "postgresql", user: "postgres",
		password: "postgresql", port: "5432",
		host:/*"db"*/ "localhost",
		database: "go-adopet",
		ssl:      "disable",
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: creds.formatStr(),
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})

	if err != nil {
		log.Fatal(err.Error())
	}
	db.Logger = logger.Default.LogMode(logger.Info)
	db.AutoMigrate()

}
