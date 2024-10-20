package db

import (
	"testing"
	// "regexp"
)

func TestReadDB(t *testing.T) {
	dbread("../_testdata/test.sqlite")
}

func TestWriteDB(t *testing.T) {
	test1("../_testdata/New.sqlite")
}
