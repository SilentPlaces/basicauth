package db

import (
	"database/sql"
	"fmt"
	"github.com/SilentPlaces/basicauth.git/internal/services/consul"
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

func NewMySQLDb(consulService *consul.ConsulService) (*sql.DB, error) {
	cfg, err := consulService.GetMySQLConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot access config: %w", err)
	}

	user := cfg[constants.MySQLUserKey]
	password := cfg[constants.MySQLPasswordKey]
	host := cfg[constants.MySQLHostKey]
	port := cfg[constants.MySQLPortKey]
	dbName := cfg[constants.MySQLDBKey]

	lifeTimeInt, err := helpers.ParseInt("max lifetime", cfg[constants.MySQLMaxLifetimeSecondsKey])
	if err != nil {
		return nil, err
	}
	maxOpenConn, err := helpers.ParseInt("max open connections", cfg[constants.MySQLMaxOpenConnectionsKey])
	if err != nil {
		return nil, err
	}
	maxIdleConn, err := helpers.ParseInt("max idle connections", cfg[constants.MySQLIdleConnectionsKey])
	if err != nil {
		return nil, err
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	db, dbErr := sql.Open("mysql", dataSourceName)
	if dbErr != nil {
		return nil, dbErr
	}

	db.SetConnMaxLifetime(time.Duration(lifeTimeInt) * time.Second)
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)

	if pingErr := db.Ping(); pingErr != nil {
		return nil, pingErr
	}
	return db, nil

}

var ProviderSet = wire.NewSet(NewMySQLDb)
