package db

import (
	"testing"
	// "regexp"
)

func TestReadDB(t *testing.T) {
	dbread("../_testdata/test.sqlite")
}
