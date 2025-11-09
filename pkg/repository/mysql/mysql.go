package mysqlrepo

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB creates a GORM MySQL connection using DSN.
// Example DSN: user:pass@tcp(127.0.0.1:3306)/airport?charset=utf8mb4&parseTime=True&loc=Local
func NewDB(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("empty DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
