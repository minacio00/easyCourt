package database

import (
	"fmt"
	"log"

	"github.com/minacio00/easyCourt/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbCreds struct {
	host     string
	user     string
	password string
	dbname   string
	port     string
	sslmode  string
}

func (cr *dbCreds) fmtString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cr.host, cr.user, cr.password, cr.dbname, cr.port, cr.sslmode,
	)
}

var Db *gorm.DB

func ConnectDb() {
	creds := &dbCreds{
		host:     config.HOST,
		user:     config.USER,
		password: config.PASSWORD,
		dbname:   config.DBNAME,
		port:     config.PORT,
		sslmode:  "disable",
	}

	Db, err := gorm.Open(postgres.Open(creds.fmtString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}

	Db.AutoMigrate(&dbCreds{})
}
