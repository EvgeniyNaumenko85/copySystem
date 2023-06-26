package db

import (
	"errors"
	"fmt"
)

var createTables = []string{
	CreateTableUsers,
	CreateTableFiles,
	CreateTableAccess,
}

func Up() error {
	for i, table := range createTables {
		_, err := GetDBConn().Exec(table)
		if err != nil {
			return errors.New(
				fmt.Sprintf("error occurred while creating table №%d, error is: %s", i, err.Error()))
		}
	}
	return nil
}
