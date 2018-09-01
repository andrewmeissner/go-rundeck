package rundeck

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// UnsafeOperations are operations that are not officially sanctioned through the API,
// yet proof useful due to API unresponsiveness.  Use with extreme caution
type UnsafeOperations struct {
	c *Client
}

// Unsafe does dastardly things
func (c *Client) Unsafe() *UnsafeOperations {
	return &UnsafeOperations{c: c}
}

type DB struct {
	u    *UnsafeOperations
	conf *DBConfig
}

type DBConfig struct {
	Dialect    string
	Host       string
	Port       int
	Username   string
	Password   string
	DBName     string
	SSLEnabled bool
}

func (u *UnsafeOperations) DB(dbConf *DBConfig) *DB {
	return &DB{u: u, conf: dbConf}
}

// ScheduledExecutions returns all scheduld executions in the database
func (db *DB) ScheduledExecutions() ([]*ScheduledExecution, error) {
	conn, err := db.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var executions []*ScheduledExecution
	return executions, conn.Find(&executions).Error
}

func (db *DB) createConnection() (*gorm.DB, error) {
	switch strings.ToLower(db.conf.Dialect) {
	case "mysql":
		return db.mysqlConnection()
	case "postgres":
		return db.postgresConnection()``
	}
	return nil, errors.New("unsupported db type")
}

func (db *DB) mysqlConnection() (*gorm.DB, error) {
	connectionStrFmt := "%s:%s@tcp(%s:%d)/%s?parseTime=true&tls=%t"
	connectionStr := fmt.Sprintf(connectionStrFmt,
		db.conf.Username,
		db.conf.Password,
		db.conf.Host,
		db.conf.Port,
		db.conf.SSLEnabled,
	)

	return gorm.Open("mysql", connectionStr)
}

func (db *DB) postgresConnection() (*gorm.DB, error) {
	connectionStrFmt := "host=%s port=%d user=%s dbname=%s password=%s sslmode=%s"
	sslMode := "disable"
	if db.conf.SSLEnabled {
		sslMode = "require"
	}

	connectionStr := fmt.Sprintf(connectionStrFmt,
		db.conf.Host,
		db.conf.Port,
		db.conf.Username,
		db.conf.DBName,
		db.conf.Password,
		sslMode,
	)

	return gorm.Open("postgres", connectionStr)
}
