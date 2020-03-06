package mssql

import (
	"os"
	"testing"
)

func TestFind(t *testing.T) {
	var columns []string
	table := os.Getenv("TABLE_TEST")
	GetDB().Table(table).Select(table).Where("column1 IS NOT NULL").Where("column2 IS NOT NULL").Find(&columns)
}
