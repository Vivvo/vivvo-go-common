package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
	"os"
	"strings"
)

// InitDB loads the config from environment variables and establishes a connection to the database
func InitDB(schema string, flags ...string) (*sql.DB, error) {
	db, err := InitDBWithoutMigrations(schema, flags...)
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
func InitDBWithoutMigrations(schema string, flags ...string) (*sql.DB, error) {
	mysqlDsn := os.Getenv("DATABASE")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbHost := os.Getenv("MYSQL_HOST")
	if dbUser != "" && dbPass != "" && dbHost != "" {
		mysqlDsn = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, schema)
		if len(flags) > 0 {
			mysqlDsn = fmt.Sprintf("%s?%s", mysqlDsn, strings.Join(flags[:], "&"))
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

	db, err := sql.Open("mysql", mysqlDsn)
	if err != nil {
		// log.Fatalf("Failed to connect to database: %s", err.Error())
		return nil, err
	}

	if err = db.Ping(); err != nil {
		// log.Fatalf("Failed to ping database: %s", err.Error())
		return nil, err
	}

	// https://stackoverflow.com/questions/39980902/golang-mysql-error-packets-go33-unexpected-eof
	db.SetMaxIdleConns(0)

	return db, nil
}

// DoMigrations from a migration folder
func DoMigrations(dir string, db *sql.DB) error {
	_ = goose.SetDialect("mysql")
	return goose.Run("up", db, dir)
}
