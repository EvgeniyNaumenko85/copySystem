package db

import (
	"errors"
	"fmt"
)

var dropTables = []string{
	DropTasksTable,
	DropUsersTable,
}

func Down() error {
	for i, table := range dropTables {
		_, err := GetDBConn().Exec(table)
		if err != nil {
			return errors.New(
				fmt.Sprintf("error occurred while dropping table №%d, error is: %s", i, err.Error()))
		}
	}
	return nil
}
