// Package gorm provides helper functions for tracing the xmirya/gorm package (https://github.com/jinzhu/gorm).
package gorm

import (
	sqltraced "github.com/meetcircle/dd-trace-go/contrib/database/sql"

	"github.com/xmirya/gorm"
)

// Open opens a new (traced) database connection. The used dialect must be formerly registered
// using (github.com/meetcircle/dd-trace-go/contrib/database/sql).Register.
func Open(dialect, source string) (*gorm.DB, error) {
	db, err := sqltraced.Open(dialect, source)
	if err != nil {
		return nil, err
	}
	return gorm.Open(dialect, db)
}
