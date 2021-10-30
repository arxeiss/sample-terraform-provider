package database

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

func handleChangeResult(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func handleInsertResult(res sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func IsUniqueConstraintError(err error) bool {
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
