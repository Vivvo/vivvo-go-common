package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
)

// InitDB loads the config from environment variables and establishes a connection to the database
func InitDB(schema string) (*sql.DB, error) {
	db, err := InitDBWithoutMigrations(schema)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat("./migrations"); err == nil {
		err = DoMigrations("./migrations", db)
		if err != nil {
			// logger.Fatalf("Failed to run migrations: %s", err.Error())
			return nil, err
		}
	}

	return db, nil
}

// InitDBWithoutMigrations if you you need to talk to a schema not owned by the service
// you shouldn't do this unless you need to, but you shouldn't need to
func InitDBWithoutMigrations(schema string) (*sql.DB, error) {
	mysqlDsn := os.Getenv("DATABASE")
	mysqlDsnSchema := fmt.Sprintf("%s%s", os.Getenv("DATABASE"), schema)
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbHost := os.Getenv("MYSQL_HOST")
	if dbUser != "" && dbPass != "" && dbHost != "" {
		mysqlDsn = fmt.Sprintf("%s:%s@tcp(%s)/", dbUser, dbPass, dbHost)
		mysqlDsnSchema = fmt.Sprintf("%s%s", mysqlDsn, schema)
		flags := os.Getenv("MYSQL_FLAGS")
		if flags != "" {
			mysqlDsn = fmt.Sprintf("%s?%s", mysqlDsn, flags)
			mysqlDsnSchema = fmt.Sprintf("%s?s%s", mysqlDsnSchema, flags)
		}
	}

	if mysqlDsn == "" {
		return nil, errors.New("no database configuration was provided")
	}

	if pem := os.Getenv("MYSQL_CERTIFICATE"); pem != "" {
		rootCertPool := x509.NewCertPool()
		if ok := rootCertPool.AppendCertsFromPEM([]byte(pem)); !ok {
			return nil, errors.New("failed to append PEM to cert pool")
		}
		if err := mysql.RegisterTLSConfig("tls_database", &tls.Config{RootCAs: rootCertPool}); err != nil {
			// log.Fatalf("Failed to load tls mysql connection: %s", err.Error())
			return nil, err
		}
	}

	db, err := sql.Open("mysql", mysqlDsnSchema)
	if err != nil {
		fmt.Printf("Failed to connect to schema, attempting to create it")
		db, err = sql.Open("mysql", mysqlDsn)
		if err != nil {
			return nil, err
		}
		_, err := db.Exec(fmt.Sprintf("create database if not exists %s", schema))
		if err != nil {
			return nil, err
		}
		db.Close()
		db, err = sql.Open("mysql", mysqlDsnSchema)
		if err != nil {
			return nil, err
		}
	}

	if err = db.Ping(); err != nil {
		// log.Fatalf("Failed to ping database: %s", err.Error())
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(2 * time.Minute)

	return db, nil
}

// DoMigrations from a migration folder
func DoMigrations(dir string, db *sql.DB) error {
	_ = goose.SetDialect("mysql")
	return goose.Run("up", db, dir)
}
