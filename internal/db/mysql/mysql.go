package mysql

import (
	"database/sql"
	"fmt"
	consul "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/convertor"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

func NewMySQLDb(consulService consul.ConsulService) (*sql.DB, error) {
	cfg, err := consulService.GetMySQLConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot access config: %w", err)
	}

	user := cfg.User
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port
	dbName := cfg.DB

	lifeTimeInt, err := helpers.ParseInt("max lifetime", cfg.MaxLifetimeSeconds)
	if err != nil {
		return nil, err
	}
	maxOpenConn, err := helpers.ParseInt("max open connections", cfg.MaxOpenConnections)
	if err != nil {
		return nil, err
	}
	maxIdleConn, err := helpers.ParseInt("max idle connections", cfg.IdleConnections)
	if err != nil {
		return nil, err
	}

	// Datasource connection string: user:password@tcp(host:port)/dbname?parseTime=true
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

var MySqlProviderSet = wire.NewSet(NewMySQLDb)
